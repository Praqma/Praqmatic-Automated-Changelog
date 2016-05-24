#This is my proposed way of adding a decorator. Use this to populate additonal attributes to the PACTask.
#We can use the key value pairs easily in liquid templates. This means that when liquid requests the 'attributes' of
#a particular task...we can fetch the data here in the decorator 
module JiraTaskDecorator
  require 'net/http' 
  require 'uri'
  require 'json'

  attr_accessor :data 

  def fetch(query_string, usr, pw, debug: false) 
    expanded = eval('"'+query_string+'"')    
	  uri = URI.parse(expanded)
     
    if expanded =~ URI::regexp
      req = Net::HTTP::Get.new(uri)
    else
      raise "Invalid URI: #{expanded}"
    end 

    #req.basic_auth usr, pw
    
    begin
      if uri.scheme == 'https'
        res = Net::HTTP.start(uri.hostname, uri.port, :use_ssl => uri.scheme = 'https') { |http|
          http.request(req)
        }
      else
        res = Net::HTTP.start(uri.hostname, uri.port) { |http|
          http.request(req)
        }        
      end
    rescue Exception => e  
      puts e.message 
      raise Exception, "Unknown host error for task with id #{task_id} on url #{expanded}"
    end 

    unless res.is_a? Net::HTTPOK
      raise Exception, "Failed to fetch task with id #{task_id} on url #{expanded} return code was #{res.code}"
    end

    begin 
      @data = parse(res.body) 
      puts "[PAC] #{@data}" if debug
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

  def fetch(debug: false)
    begin
      unless task_id.nil? 
        ticket = TracTaskDecorator.trac_instance.tickets.get task_id.to_i
        @data = { :summary => ticket.summary, :status => ticket.status, :description => ticket.description }
        puts "[PAC] #{@data}" if debug
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
