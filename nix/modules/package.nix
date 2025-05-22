{ lib, buildGoModule }:

buildGoModule {
  pname = "cablebox-control";
  version = "0.1.0";
  src = ../../..;
  meta = with lib; {
    description = "Cablebox control service";
    homepage = "https://github.com/scottjab/cablebox-control";
    license = licenses.mit;
    maintainers = with maintainers; [ scottjab ];
    platforms = platforms.unix;
  };
} 