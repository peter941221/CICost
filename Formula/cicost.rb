class Cicost < Formula
  desc "GitHub Actions cost and waste hotspot analyzer"
  homepage "https://github.com/peter941221/CICost"
  version "0.1.0"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/peter941221/CICost/releases/download/v#{version}/cicost_#{version}_darwin_arm64.tar.gz"
      sha256 "REPLACE_WITH_REAL_SHA256"
    else
      url "https://github.com/peter941221/CICost/releases/download/v#{version}/cicost_#{version}_darwin_amd64.tar.gz"
      sha256 "REPLACE_WITH_REAL_SHA256"
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/peter941221/CICost/releases/download/v#{version}/cicost_#{version}_linux_arm64.tar.gz"
      sha256 "REPLACE_WITH_REAL_SHA256"
    else
      url "https://github.com/peter941221/CICost/releases/download/v#{version}/cicost_#{version}_linux_amd64.tar.gz"
      sha256 "REPLACE_WITH_REAL_SHA256"
    end
  end

  def install
    bin.install "cicost"
  end

  test do
    assert_match "cicost", shell_output("#{bin}/cicost version")
  end
end

