{
  lib,
  buildGoModule,
  stdenv,
}:

buildGoModule rec {
  pname = "cablebox-control";
  version = "0.1.0";

  src = ./.;

  vendorHash = null; # This will be filled in automatically on first build

  # Build the main package
  buildPhase = ''
    export GOPATH=$TMPDIR/go
    mkdir -p $GOPATH/src/github.com/scottjab
    cp -r $src $GOPATH/src/github.com/scottjab/cablebox-control
    cd $GOPATH/src/github.com/scottjab/cablebox-control
    go build -o $out/bin/cablebox-control ./cmd/cablebox-control
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
