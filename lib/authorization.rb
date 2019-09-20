module Authorization

  def self.create(config = {}, uri = nil)
    if config.has_key?(:github)
      Authorization::GithubToken.new(config[:github])
    elsif config.has_key?(:basic)
      Authorization::BasicAuth.new(config[:basic], uri)
    end
  end

  class GithubToken
    def initialize(config = {})
      @config = config
    end

    def headers
      heads = { 'Authorization' => "token " + @config[:token], }
      heads.merge(@config[:headers] || {})
      heads
    end

    def to_s
      "GitHub token authentication"
    end
  end

  class BasicAuth
    def initialize(config = {}, uri)
      @uri = uri
      @config = config
    end

    def headers
      if @config.has_key?(:netrc)
        if @config[:netrc].nil?
          netrc_config = Netrc.read
        else
          netrc_config = Netrc.read(@config[:netrc])
        end
        netrc_usr, netrc_pw = netrc_config[URI.parse(@uri).host]
        heads = { 'Authorization' => "Basic " + Base64.encode64(netrc_usr + ":" + netrc_pw), }
      else
        heads = { 'Authorization' => "Basic " + Base64.encode64(@config[:username] + ":" + @config[:password]), }
      end
      heads.merge!(@config[:headers] || {})
      heads
    end

    def to_s
      "Basic authentication"
    end
  end
end
