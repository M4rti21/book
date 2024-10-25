pkgname=book-bin
pkgver=0.2
pkgrel=1
pkgdesc="Simple bookmark manager written in go"
arch=("x86_64")
url="https://github.com/m4rti21/book"
liscence="GPL-3.0"
source=("https://github.com/m4rti21/${pkgname}/releases/download/${pkgver}/${pkgname}") # Adjust the URL for the binary
noextract=("${pkgname}")
sha256sums=('SKIP')
options=('!debug')

package() {
    install -Dm755 "${pkgname}" "$pkgdir/usr/bin/${pkgname}"
}
