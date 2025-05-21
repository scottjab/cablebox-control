{ lib, buildGoModule }:

buildGoModule {
  pname = "cablebox-control";
  version = "0.1.0";
  src = ../../..;
  vendorSha256 = "sha256-0000000000000000000000000000000000000000000000000000";
  meta = with lib; {
    description = "Cablebox control service";
    homepage = "https://github.com/scottjab/cablebox-control";
    license = licenses.mit;
    maintainers = with maintainers; [ scottjab ];
    platforms = platforms.unix;
  };
} 