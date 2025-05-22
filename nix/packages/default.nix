{
  lib,
  buildGoModule,
  stdenv,
  src,
}:

buildGoModule rec {
  pname = "cablebox-control";
  version = "0.1.0";

  inherit src;

  vendorHash = null; # This will be filled in automatically on first build

  # Specify the path to the main package relative to the root
  subPackages = [ "cmd/cablebox-control" ];


  meta = with lib; {
    description = "Cablebox control application";
    homepage = "https://github.com/scottjab/cablebox-control";
    license = licenses.mit;
    maintainers = with maintainers; [
      "scottjab@gmail.com"
    ];
    platforms = platforms.unix;
  };
}
