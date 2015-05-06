# https://github.com/gemhome/rmagick/pull/132/files with some changes
if [ "$TRAVIS_OS_NAME" == "osx" ]; then
  brew update
  brew install imagemagick
  brew install webp
  exit
fi

if [ "$TRAVIS_OS_NAME" == "linux" ]; then
  IMAGEMAGICK_VERSION=6.9.1-1

  wget http://downloads.webmproject.org/releases/webp/libwebp-0.4.3.tar.gz
  tar -xzf libwebp-0.4.3.tar.gz
  cd libwebp-0.4.3
  mkdir tmp
  ./configure --disable-shared --prefix=tmp --bin-dir=tmp
  make && make install
  cd ..

  dpkg --list imagemagick
  sudo apt-get remove imagemagick
  sudo apt-get install build-essential libx11-dev libxext-dev zlib1g-dev libpng12-dev libjpeg-dev libfreetype6-dev libxml2-dev
  sudo apt-get build-dep imagemagick

  wget http://www.imagemagick.org/download/releases/ImageMagick-${IMAGEMAGICK_VERSION}.tar.gz
  tar -xzf ImageMagick-${IMAGEMAGICK_VERSION}.tar.gz
  cd ImageMagick-${IMAGEMAGICK_VERSION}
  ./configure --prefix=/usr
  sudo make install
  cd ..
  sudo ldconfig
  exit
fi
