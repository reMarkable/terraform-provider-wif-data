{
  pkgs,
  ...
}:

{
  pre-commit = {
    hooks = {
      gofmt.enable = true;
      shellcheck.enable = true;
      markdownlint = {
        enable = true;
      };
      yamllint.enable = true;
      commitizen.enable = true;
    };
  };
  # https://devenv.sh/packages/
  packages = [ pkgs.go-task ];
  languages.go.enable = true;
  enterTest = ''
    echo "Entering test environment"
    task test
  '';
}
