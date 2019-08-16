#This is my proposed way of adding a decorator. Use this to populate additonal attributes to the PACTask.
#We can use the key value pairs easily in liquid templates. This means that when liquid requests the 'attributes' of
#a particular task...we can fetch the data here in the decorator
module JiraTaskDecorator
  require 'net/http'
  require 'uri'
  require 'json'
  require 'base64'
  require_relative 'logging'

  attr_accessor :data

  def fetch(query_string, usr, pw)
    expanded = eval('"'+query_string+'"')
	  uri = URI.parse(expanded)

    begin
      res = DecoratorUtils.query(uri, usr, pw)
    rescue Exception => net_error
      raise Exception, "Unknown host error for task with id #{task_id} on url #{expanded}\n#{net_error}"
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

module TracTaskDecorator
  require 'trac4r'

  def self.trac_instance
    @@trac_instance
  end

  def self.trac_instance=(trac)
    @@trac_instance = trac
  end

  def fetch
    begin
      unless task_id.nil?
        ticket = TracTaskDecorator.trac_instance.tickets.get task_id.to_i
        @data = { :summary => ticket.summary, :status => ticket.status, :description => ticket.description }
        Logging.verboseprint(1, "[PAC] Fetched the following from Trac: #{@data}")
        @data
      end
    rescue Trac::TracException => e
      raise Exception, "[PAC] The ticket with the id #{task_id} not found in Trac"
    end
  end

  def attributes
    super.merge!(
      {
        'data' => @data
      }
    )
  end

  attr_accessor :data
end

module FogbugzTaskDecorator

  require 'xmlsimple'

  def fetch(query_string)
    expanded = eval('"'+query_string+'"')
    uri = URI.parse(expanded)
    res = DecoratorUtils.query(uri)
    Logging.verboseprint(3, "[PAC] Got the following data from #{expanded}: #{res.body}")
    @data = XmlSimple.xml_in res.body
    Logging.verboseprint(1, "[PAC] Fetched the following from FogBugz: #{@data}")
    @data
  end

  def attributes
    super.merge!({ 'data' => @data })
  end
end

module DecoratorUtils extend self

  def query(uri, usr = nil, pw = nil)
    req = Net::HTTP::Get.new(uri)
    unless usr.nil?
      Logging.verboseprint(3, "[PAC] Using basic authentication")
      req['Authorization'] = "Basic " + Base64.encode64(usr+":"+pw)
      req['Content-Type'] = "application/json"
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
