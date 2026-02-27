class Cicost < Formula
  desc "GitHub Actions cost and waste hotspot analyzer"
  homepage "https://github.com/peter941221/CICost"
  version "0.2.0"
  license "Apache-2.0"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/peter941221/CICost/releases/download/v#{version}/cicost_#{version}_darwin_arm64.tar.gz"
      sha256 "98ad44cf993365f9108bc6aec84307651fbd085170a67cec791547a56e166884"
    else
      url "https://github.com/peter941221/CICost/releases/download/v#{version}/cicost_#{version}_darwin_amd64.tar.gz"
      sha256 "5affd94117181f772299f5496f5d03e8d7538356970d897e9a311181be28b426"
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/peter941221/CICost/releases/download/v#{version}/cicost_#{version}_linux_arm64.tar.gz"
      sha256 "64efc35f40a0896199a38a385bacb2937195af44dce4c2c8b549d1b37473fa54"
    else
      url "https://github.com/peter941221/CICost/releases/download/v#{version}/cicost_#{version}_linux_amd64.tar.gz"
      sha256 "9c28146741e4a5bd1624f37e088dd54ac81db8fe140e24c082126730540097aa"
    end
  end

  def install
    bin.install "cicost"
  end

  test do
    assert_match "cicost", shell_output("#{bin}/cicost version")
  end
end
