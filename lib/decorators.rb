#This is my proposed way of adding a decorator. Use this to populate additonal attributes to the PACTask.
#We can use the key value pairs easily in liquid templates. This means that when liquid requests the 'attributes' of
#a particular task...we can fetch the data here in the decorator 
module JiraTaskDecorator
  require 'net/http' 
  require 'uri'
  require 'json'

  def fetch(query_string, usr, pw) 
    expanded = eval('"'+query_string+'"')

    uri = URI(expanded)
    req = Net::HTTP::Get.new(uri)
    req.basic_auth usr, pw
    res = Net::HTTP.start(uri.hostname, uri.port) { |http|
      http.request(req)
    }

    @data = JSON.parse(res.body)
  end

  def attributes
    super.merge!(
      { 
        :data => @data
      }
    )
  end

  attr_accessor :data  
end

module CrucibleTaskDecorator
  def attributes
    super.merge! ({ :review_status => 'APPROVED' })
  end 
end

module TracTaskDecorator
  require 'trac4r'

  def trac_instance
    @@trac_instance    
  end

  def trac_instance=(trac)
    @@trac_instance = trac
  end

  def fetch
    begin
      unless task_id.nil? 
        ticket = track_instance.tickets.get array[task_id.to_i]
        @data = { :summary => ticket.summary, :status => ticket.status, :description => ticket.description }
      end
    rescue Trac::TracException => e
      puts "The ticket with the id #{task_id} not found in Trac"
      puts e.message
    end    
  end

  def attributes
    super.merge!(
      { 
        :data => @data
      }
    )    
  end

  attr_accessor :data  
end