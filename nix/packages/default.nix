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

  subPackages = [ "cmd/cablebox-control" ];

  preBuild = ''
    cd $sourceRoot
  '';

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
