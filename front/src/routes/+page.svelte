<script>
  import { onMount, onDestroy } from 'svelte';

  let messages = [];
  let inputMessage = "";
  let ws;

  onMount(() => {
    // 1. GleamのサーバーにWebSocket接続（裏側でGETリクエストが飛ぶ）
    ws = new WebSocket("ws://localhost:8000/chat");

    // 2. Gleamからメッセージが送られてきた時の処理
    ws.onmessage = (event) => {
      messages = [...messages, event.data];
    };
  });

  onDestroy(() => {
    if (ws) ws.close();
  });

  function sendMessage() {
    if (ws && inputMessage) {
      ws.send(inputMessage); // Gleamへ送信
      inputMessage = "";
    }
  }
</script>

<div class="chat-container">
  <div class="messages">
    {#each messages as msg}
      <p>{msg}</p>
    {/each}
  </div>
  
  <input type="text" style="color: black;"  bind:value={inputMessage} placeholder="メッセージを入力..." />
  <button on:click={sendMessage}>送信</button>
</div>
