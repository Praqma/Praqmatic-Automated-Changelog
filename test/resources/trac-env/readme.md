The Dockerfile originates from:

 * https://github.com/jmmills/docker-trac/
 * SHA: 6c751b130fe43e436d4a7515106a6f07c63b94a1

We modified it to fit our purposes using Trac to test PAC.

Changes done:

 * Switched baseimage from ubuntu:quantal to ubutu:14.04. Quantal has been deprecated
 * Removed `RUN apt-get install -y trac-batchmodify`, batchmodify is included in Trac 1.X which comes with ubutu 14.04
 * Set `TRAC_PASS` to a fixed value in `setup_trac_config.sh`. 
 * Enabled `tracrpc` in `setup_trac_config.sh` with this line: `trac-admin /trac config set components tracrpc.* enabled`
 * Added the necessary permissions: `trac-admin /trac permission add anonymous XML_RPC` and `trac-admin /trac permission add anonymous TICKET_CREATE` in `setup_trac_config.sh`
