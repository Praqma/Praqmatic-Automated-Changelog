#This is my proposed way of adding a decorator. Use this to populate additonal attributes to the PACTask.
#We can use the key value pairs easily in liquid templates. This means that when liquid requests the 'attributes' of
#a particular task...we can fetch the data here in the decorator
module JsonTaskDecorator
  require 'net/http'
  require 'uri'
  require 'json'
  require 'base64'
  require 'openssl'
  require_relative 'logging'

  attr_accessor :data

  def fetch(query_string, auth, ssl_verify)
    expanded = eval('"'+query_string+'"')
	  uri = URI.parse(expanded)

    begin
      res = DecoratorUtils.query(uri, auth, ssl_verify)
    rescue Exception => net_error
      raise Exception, "Unknown host error for task with id #{task_id} on url #{expanded}\n#{net_error}"
    end

    unless res.is_a? Net::HTTPOK
      raise Exception, "Failed to fetch task with id #{task_id} on url #{expanded} return code was #{res.code}"
    end

    begin
      Logging.verboseprint(3, "[PAC] Got the following data from #{expanded}: #{res.body}")
      @data = parse(res.body)
      Logging.verboseprint(2, "[PAC] Fetched the following json data: #{@data}")
      @data
    rescue JSONError
      raise Exception, "Unparsable JSON data fetched from url #{expanded}"
    end

  end

  def parse(response)
    JSON.parse(response)
  end

  def attributes
    super.merge!(
      {
        'data' => @data
      }
    )
  end

end

module DecoratorUtils extend self
  require_relative 'authorization'

  def query(uri, auth, ssl_verify)
    req = Net::HTTP::Get.new(uri)
    verification = ssl_verify ? OpenSSL::SSL::VERIFY_PEER : OpenSSL::SSL::VERIFY_NONE
    auth = Authorization.create(auth)
    Logging.verboseprint(1, "[PAC] Verification (0 is off, 1 is peer authentication): #{verification}")
    if auth.nil?
      Logging.verboseprint(0, "[PAC] Using no authentication")
    else
      Logging.verboseprint(0, "[PAC] Using #{auth}")
      req['Content-Type'] = "application/json"
      auth.headers.each { |k,v|
        req[k] = v
      }
    end
    res = Net::HTTP.start(uri.hostname, uri.port, :use_ssl => uri.scheme == 'https', :verify_mode => verification ) { |http|
      http.request(req)
    }
  end

end
