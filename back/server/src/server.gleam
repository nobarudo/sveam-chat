import gleam/dict
import gleam/erlang/process
import gleam/option.{Some, None}
import gleam/otp/actor
import mist
import wisp
import wisp/wisp_mist
import gleam/io
import gleam/int
import gleam/http/request
import gleam/json
import gleam/http
import gleam/httpc

@external(erlang, "logger_ffi", "add_file_handler")
fn add_file_handler(filename: String) -> Nil

pub type HubMessage {
  AddClient(id: String, subject: process.Subject(String))
  RemoveClient(id: String)
  Broadcast(id: String, text: String)
}

fn hub_loop(clients: dict.Dict(String, process.Subject(String)), message: HubMessage) {
  case message {
    AddClient(id, subject) -> {
      actor.continue(dict.insert(clients, id, subject))
    }
    RemoveClient(id) -> {
      actor.continue(dict.delete(clients, id))
    }
    Broadcast(id, text) -> {
      wisp.log_info("Message broadcasted from " <> id)
      let msg = id <> ": " <> text
      dict.each(clients, fn(_id, subject) {
        process.send(subject, msg)
      })
      actor.continue(clients)
    }
  }
}

pub fn main() {
  wisp.configure_logger()
  add_file_handler("server.log")
  let secret_key_base = wisp.random_string(64)

  let assert Ok(actor.Started(_pid, hub_subject)) = 
    actor.new(dict.new())
    |> actor.on_message(hub_loop)
    |> actor.start

  let router = fn(req: request.Request(mist.Connection)) {
    case request.path_segments(req) {
      ["chat"] -> {
        mist.websocket(
          request: req,
          on_init: fn(_conn) {
            let client_id = "User-" <> wisp.random_string(6)
           
            wisp.log_info("New client connected: " <> client_id)

            let client_subject = process.new_subject()
            let selector = 
              process.new_selector()
              |> process.select(client_subject)
            
            process.send(hub_subject, AddClient(client_id, client_subject))
            process.send(hub_subject, Broadcast("System", client_id <> " が入室しました！"))

            #(#(client_id, hub_subject), Some(selector))
          },          on_close: fn(state) {
            let #(client_id, hub) = state
            process.send(hub, RemoveClient(client_id))
            process.send(hub, Broadcast("System", client_id <> " が退出しました。"))
            Nil
          },
          handler: handle_ws_message
        )
      }
      _ -> wisp_mist.handler(handle_http_request, secret_key_base)(req)
    }
  }

  let assert Ok(_) = router |> mist.new |> mist.port(8000) |> mist.start
  process.sleep_forever()
}

fn handle_ws_message(state: #(String, process.Subject(HubMessage)), message, conn) {
  let #(client_id, hub_subject) = state

  case message {
    mist.Text(text) -> {
      let body = 
        json.object([
          #("client_id", json.string(client_id)),
          #("message", json.string(text)),
        ])
        |> json.to_string

      // 2. POSTリクエストの構築 (URLをGo側のAPIエンドポイントに合わせます)
      let assert Ok(req) = request.to("http://localhost:8080/")
      let req = 
        req
        |> request.set_method(http.Post)
        |> request.set_body(body)
        |> request.set_header("content-type", "application/json")

      // 3. 送信
      let assert Ok(resp) = httpc.send(req)

      io.println("✅ POST送信完了！")
      io.println("ステータスコード: " <> int.to_string(resp.status))
      io.println("レスポンスボディ: " <> resp.body)

      process.send(hub_subject, Broadcast(client_id, text))
      mist.continue(state)
    }
    mist.Custom(msg) -> {
      let assert Ok(_) = mist.send_text_frame(conn, msg)
      mist.continue(state)
    }
    mist.Binary(_) -> mist.continue(state)
    mist.Closed | mist.Shutdown -> mist.stop()
  }
}

fn handle_http_request(_req: wisp.Request) -> wisp.Response {
  wisp.html_response("<h1>API Server</h1>", 200)
}
