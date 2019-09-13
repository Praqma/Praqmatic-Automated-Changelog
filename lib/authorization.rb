module Authorization

  def self.create(config = {})
    if config.has_key?(:github)
      Authorization::GithubToken.new(config[:github])
    elsif config.has_key?(:basic)
      Authorization::BasicAuth.new(config[:basic])
    end
  end

  class GithubToken
    def initialize(config = {})
      @config = config
    end

    def headers
      heads = { 'Authorization' => "token " + eval(@config[:token]), }
      puts heads
      heads
    end

    def to_s
      "GitHub token authentication"
    end
  end

  class BasicAuth
    def initialize(config = {})
      @config = config
    end

    def headers
      heads = { 'Authorization' => "Basic " + Base64.encode64(eval(@config[:username]) + ":" + eval(@config[:basic][:password])), }
      puts heads
      heads
    end

    def to_s
      "Basic authentication"
    end
  end
end
