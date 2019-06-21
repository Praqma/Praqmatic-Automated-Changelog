#This is my proposed way of adding a decorator. Use this to populate additonal attributes to the PACTask.
#We can use the key value pairs easily in liquid templates. This means that when liquid requests the 'attributes' of
#a particular task...we can fetch the data here in the decorator
module JsonTaskDecorator
  require 'net/http'
  require 'uri'
  require 'json'
  require_relative 'logging'

  attr_accessor :data

  def fetch(query_string, usr, pw)
    expanded = eval('"'+query_string+'"')
	  uri = URI.parse(expanded)

    begin
      res = DecoratorUtils.query(uri, usr, pw)
    rescue Exception
      raise Exception, "Unknown host error for task with id #{task_id} on url #{expanded}"
    end

    unless res.is_a? Net::HTTPOK
      raise Exception, "Failed to fetch task with id #{task_id} on url #{expanded} return code was #{res.code}"
    end

    begin
      Logging.verboseprint(3, "[PAC] Got the following data from #{expanded}: #{res.body}")
      @data = parse(res.body)
      Logging.verboseprint(1, "[PAC] Fetched the following from Jira: #{@data}")
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

  def query(uri, usr = nil, pw = nil)
    req = Net::HTTP::Get.new(uri)
    unless usr.nil?
      Logging.verboseprint(3, "[PAC] Using basic authentication")
      req.basic_auth usr, pw
    end

    if uri.scheme == 'https'
      res = Net::HTTP.start(uri.hostname, uri.port, :use_ssl => uri.scheme = 'https') { |http|
        http.request(req)
      }
    else
      res = Net::HTTP.start(uri.hostname, uri.port) { |http|
        http.request(req)
      }
    end
  end

end
