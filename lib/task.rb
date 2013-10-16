# encoding: utf-8
require 'trac4r'
require 'kramdown'
require 'pdfkit' 

begin 
  require 'fogbugz'
rescue LoadError => error
  
end

class String 
  def is_number? 
    !!(self =~ /^[-+]?[0-9]+$/)
  end
end
  

module Task  

  #The task system is responsible for writing the changelog. We feed it with a list of commits, and an output directory for the changelogs
  #unless otherwise specified the path will be the current directory.
  class TaskSystem
    attr_accessor :path, :settings
    
    def initialize(settings, path = nil)
      @settings = settings
      @path = path
    end
    
    def footer
      "#{@settings[:general]['changelog_name']} generated #{Time.new}"
    end
    
    def footer_html
      "<span style='font-style:italic;font-size:small'>#{footer}</span>"
    end 
    
    def header_html
      "<title>#{@settings[:general]['changelog_name']}</title>"
    end
    
    def html_escape_non_ascii(text) 
      text.gsub(/Æ/,'&AElig;').gsub(/æ/,'&aelig;').gsub(/Ø/,'&Oslash;').gsub(/ø/,'&oslash;').gsub(/Å/,'&Aring;').gsub(/å/,'&aring;')
    end
  end
    
  class FogBugzTaskSystem < TaskSystem
    
    attr_accessor :instance, :fields
    SCM_CASE_REGEX = /[Case|\[Case\]|fixed]\s(?<id>([0-9]+))+/i
    SCM_JENKINS_RELATED_COMMIT = /(?<id>JENKINS-[0-9]+)/i
    
    def initialize(settings)
      super(settings)
      @instance = Fogbugz::Interface.new( :email => @settings[:fogbugz]["fogbugz_usr"], :password => @settings[:fogbugz]["fogbugz_pwd"], :uri => @settings[:fogbugz]["fogbugz_url"] )
      @instance.authenticate
    end
    
    #Filter commits method. This is wehre we apply our regexes.
    def filter(cmits)
      caseNumbers  = []
      titleSearchStrings = []

      cmits.each do |cMsg|
        
        res = Regexp.union(SCM_CASE_REGEX).match(cMsg)       
        if !res.nil?
          if res[:id].is_number?            
            caseNumbers << res[:id]
          else 
            titleSearchStrings << res[:id]
          end
        end                
      end
      
      #If the results of the parsing yields something that is not a case number (IE 1200, 2332 etc). Each item has to seperated with OR.
      if titleSearchStrings.length > 0         
        xmloutput_title = @instance.command(:search, :q => "title:"+titleSearchStrings.join(" OR "), :cols => @settings[:fogbugz]["fogbugz_fields"])            
      end
      
      cases_xml_array = []
      if caseNumbers.length > 0
        xmloutput_cases = @instance.command(:search, :q => caseNumbers.uniq.join(","), :cols => @settings[:fogbugz]["fogbugz_fields"])        
        caseNumbers.each do |c| 
          case_xml = @instance.command(:search, :q => c, :cols => @settings[:fogbugz]["fogbugz_fields"])
          cases_xml_array.push(case_xml)
        end        
      end
      
      hashes = []
      unless cases_xml_array.empty?
        cases_xml_array.each do |x|
          unless x.nil?
             h = Hash.new
              x["cases"]["case"].each do |aKey, aValue|
                h[aKey.to_sym] = aValue
              end
              hashes.push(h)           
          end
        end
      end
      descending = -1
      hashes.uniq { |x| x[:ixBug] }.sort_by { |hv| hv[:ixBug].to_i * descending  }
    end

    #Markup version of the task link.
    def case_link_md(caseNumber)
      "[Case #{caseNumber}](#{@settings[:fogbugz]["fogbugz_url]"]}/default.asp?#{caseNumber})"
    end
    
    #Html version of the case link
    def case_link_html(caseNumber) 
      "<a class='task-link' href='#{@settings[:fogbugz]['fogbugz_url']}/default.asp?#{caseNumber}'>#{caseNumber}</a>"
    end
    
    #By default we use the css in here. If a css file is specified, then that is used instead.
    def css_style
      if @settings[:general]['changelog_css'].nil?
        return  "<link rel='stylesheet' type='text/css' href='#{File.join(File.dirname(__FILE__),'..','default_changelog_css.css')}'/>" 
      else
        return "<link rel='stylesheet' type='text/css' href='#{@settings[:general]['changelog_css']}'/>" 
      end
    end
    
    #Returns a value representing the percentage of referenced commits
    def health(commits)
      taskReferenceCount = commits.select { |e| !Regexp.union(SCM_CASE_REGEX).match(e).nil? }.length
      unspecifiedCommitCount = commits.length - taskReferenceCount
      return (taskReferenceCount / commits.length.to_f)       
    end   
    
    def statistics_html(commits)
      taskReferenceCount = commits.select { |e| !Regexp.union(SCM_CASE_REGEX).match(e).nil? }.length
      unspecifiedCommitCount = commits.length - taskReferenceCount
      healthgauge = health(commits)*100
      html =  <<-EOM
        <div id="metadata-section"> 
          <h2 id="metadata-details-header">Details</h2>
          <div id='metadata'>
          <p>            
            This changelog contains<strong> #{commits.length}</strong> commits<br/>
            Number of referenced commits is <strong>#{taskReferenceCount}</strong> which is <strong>#{healthgauge}%</strong> of all commits<br/>
            Which leaves out <strong>#{unspecifiedCommitCount}</strong> commits without proper commit messages<br/>            
            #{footer_html}
          </p>
          </div>
        </div>
      EOM
     
    end
    
    def write_markdown(commits)             
      #Use a file block instead. This way we avoid having to expricitly close the file
      mdPath = @path.nil? == true ? "changelog.md" : File.join(@path, "#{@settings[:general]['changelog_name']}.md")
      result = filter(commits)                 
      File.open(mdPath.to_s, 'w') do |file|
        file << "# #{@settings[:general]['changelog_name']}\n"
        file << "\n"
        unless result.nil?
          result.each do |task|
            unless task[:sStatus].downcase.include? "duplicate"
              file << "\n"
              
              unless task[:sReleaseNotes].nil?
                file << "**Release note:** #{task[:sReleaseNotes]}\n"
                file << "\n"
              end
              
              file << "**Status:** #{task[:sStatus]}\n"
              file << "\n"
              file << "**Title:** #{task[:sTitle]}\n"
              file << "\n"
            end                    
          end
        end
        
        file << footer        
      end
    end
     
    #Write HTML uses the release-notes field from fogbugz to generate a a list of release notes. 
    def write_html(commits)
      htmlPath = @path.nil? == true  ? "changelog.html" : File.join(@path, "#{@settings[:general]['changelog_name']}.html")
      tasks = filter(commits)
      
      File.open(htmlPath,'w+') do |file|
        file << "<html>"
        file << "<head>"        
        file << css_style
        file << header_html
        file << "</head>"
        file << "<body>"
        file << "<div id='changelog'>"
        file << "<h1 id='changelog-headline'>#{@settings[:general]['changelog_name']}</h1>"
        file << "<div id='change-list'>"
        unless tasks.nil?
          tasks.each do |task|
            puts task[:sStatus]
            unless task[:sStatus].downcase.include? "duplicate"              
            file << "<div class ='change-item'>"
            file << ""
            
            title_case = "<div class='change-title' >
            <div class='issue-type-#{task[:sCategory]} issue-type'>&nbsp;</div>
            #{task[:sTitle]} (#{case_link_html(task[:ixBug])})</div>"
            title_case_escaped = html_escape_non_ascii(title_case)
            file << title_case_escaped
             
             unless task[:sReleaseNotes].nil?
                rel_note = "<div class='release-note'>#{task[:sReleaseNotes]} </div>"
                rel_note_escaped = html_escape_non_ascii(rel_note)
                file << rel_note_escaped
             end
             file << '</div>'                        
            end
          end
        end
        file << "</div>"
        file << statistics_html(commits)
        file << "</div>"
        file << "</body>"
        file << "</html>"
      end  
    end
    
    
    def write_changelog(commits, path)
      @path = path
      write_markdown(commits)
      unless @settings[:general]['changelog_formats'].nil?
        if @settings[:general]['changelog_formats'].include?("html")
          write_html(commits)
        end
      end       
    end
            
  end
  
  class TracTaskSystem < TaskSystem
    TASK_REGEX = /Ticket\#(?<id>([0-9]+|none))+/i
    attr_accessor :instance, :activePattern
    
    def initialize(settings)
      super(settings) 
      @instance = Trac.new @settings[:trac]['trac_url'], @settings[:trac]['trac_usr'], @settings[:trac]['trac_pwd']
    end
    
    def write_markdown(commits, path)
      md_path = path.nil? == true ? "changelog.md" : File.join(path, "#{@settings[:general]['changelog_name']}.md")
      File.open(md_path,'w:UTF-8') do |file|
        file << "# #{@settings[:general]['changelog_name']}\n"
        file << "\n"
        
        taskReferenceCount = commits.select { |e| !e.match(TASK_REGEX).nil? }.length 
        noneTaskCount = commits.select { |e| !e.match(TASK_REGEX).nil? && e.match(TASK_REGEX)[:id].casecmp("none") == 0 }.length
        unspecifiedCommitCount = commits.length - taskReferenceCount

        file << "This changelog encompasses **#{commits.length}** commits\n"
        file << "\n"
        file << "**#{taskReferenceCount}** Task references found"
        file << "\n"
        file << "**#{noneTaskCount}** 'none' task reference found"
        file << "\n\n"
        
        filtered_data = filter(commits)

        filtered_data.each do |task|
          summary_line = "### (Ticket #{task[:id]}) #{task[:summary]} \n\n"          
          file << summary_line         
          
          if !task[:status].nil?            
            status_line = "**Status:** #{task[:status]}\n\n"
            file << status_line
          end
          
          if !task[:description].nil?
            description_line = "**Description:** #{task[:description]}\n\n"
            file << description_line            
          end    
        end
        file << footer
 
      end
      md_path  
    end
    
    def write_html(file_path_to_markdown)
      input = File.read(file_path_to_markdown)
      html = Kramdown::Document.new(input, :smart_quotes => 'apos,apos,quot,quot', :external_encoding => 'UTF-8').to_html
      output_file_path = file_path_to_markdown.sub(/\.md$/, '.html')
      File.open(output_file_path,'w:utf-8') do |file|        
        file << html.gsub(/[æ]/,'&aelig;').gsub(/[Æ]/,'&AElig;').gsub(/ø/,'&oslash;').gsub(/Ø/,'&Oslash;').gsub(/Å/,'&Aring;').gsub(/å/, '&aring;').gsub(/”/,'"')
      end
    end
    
    def write_pdf(file_path_to_markdown)
      html = Kramdown::Document.new(File.read(file_path_to_markdown)).to_html
      output_file_path = file_path_to_markdown.sub(/\.md$/, '.pdf')
      kit = PDFKit.new(html, :page_size => 'A4')
      kit.to_file(output_file_path) 
    end
    
    def write_changelog(commits, path)
      file_path_markdown = write_markdown(commits, path)
      if @settings[:general]['changelog_formats'].include?("html")
        write_html(file_path_markdown)
      end
      if @settings[:general]['changelog_formats'].include?("pdf")
        write_pdf(file_path_markdown)
      end
    end
    
    def filter(commit)
      references = []
      commit.each do |co| 
        match = co.match(TASK_REGEX)
        if !match.nil?
          hash = { :id => match[:id] }
          references.push(hash)
        end 
      end
          
      references.sort { |x,y| y[:id] <=> x[:id] }
      
      references.each do |array|
        begin
          if array[:id].casecmp("none") != 0      
            ticket = instance.tickets.get array[:id]
            array[:summary] = ticket.summary
            array[:status] = ticket.status
            array[:description] = ticket.description
          end        
        rescue Trac::TracException => e                 
          puts "The ticket with the id #{array[:id]} not found in Trac"
          puts e.message          
        end
      end

      filteredList = references.uniq { |x| x[:id] }      
      filteredList = filteredList.reject { |e| e[:summary].nil? }                   
    end   
  end  
end