{
  description = "sveam-chat: Gleam & SvelteKit Development Environment";

  inputs = {
    # 常に最新の安定版ツールチェインを取得する
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
  };

  outputs = { self, nixpkgs }:
    let
      # サポートするOSアーキテクチャのリスト
      supportedSystems = [ "x86_64-linux" "aarch64-linux" "x86_64-darwin" "aarch64-darwin" ];
      forEachSystem = nixpkgs.lib.genAttrs supportedSystems;
    in
    {
      devShells = forEachSystem (system:
        let
          pkgs = import nixpkgs { inherit system; };
        in
        {
          default = pkgs.mkShell {
            # 🌟 sveam-chatの開発に必要なパッケージ群
            buildInputs = with pkgs; [
              # --- Backend (Gleam) ---
              gleam       # Gleamコンパイラ
              erlang_27   # Gleamが裏で動くためのErlang VM
              rebar3      # Erlang系の依存解決ビルドツール（Gleamが内部で呼び出します）

              # --- Frontend (SvelteKit) ---
              nodejs_22   # Node.js環境
              # pnpm や yarn を使っている場合はここに追加（例: nodePackages.pnpm）
            ];

            # 🌟 シェルが立ち上がった時に自動実行されるスクリプト
            shellHook = ''
              echo "========================================="
              echo " ⚡ sveam-chat Nix Environment Active"
              echo "========================================="
              echo "🟢 $(gleam --version)"
              echo "🟢 Node.js $(node --version)"
              echo "========================================="
            '';
          };
        });
    };
}
