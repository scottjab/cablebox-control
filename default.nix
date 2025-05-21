{ lib
, buildGoModule
, fetchFromGitHub
}:

buildGoModule rec {
  pname = "cablebox-control";
  version = "0.1.0";

  src = ./.;

  vendorHash = null; # This will be filled in automatically on first build

  meta = with lib; {
    description = "Cablebox control application";
    homepage = "https://github.com/scottjab/cablebox-control";
    license = licenses.mit;
    maintainers = with maintainers; [ 
      "scottjab@gmail.com"
    ];
  };
} 