# encoding: utf-8
require 'kramdown'
require 'pdfkit'
begin
  require 'trac4r'
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
    attr_accessor :path, :settings, :filterpatterns
    def initialize(settings, path = nil)
      @settings = settings
      @path = path
      @filterpatterns = filterpatterns
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

    def statistics_html(commits)
      tasks = task_id_list(commits)
      taskReferenceCount = commits.values.select { |e| !@filterpatterns.match(e).nil? }.length
      unspecifiedCommitCount = get_shas_without_reference(commits, tasks).length
      healthgauge = health(commits)*100
      html =  
      <<-EOM
<div id="metadata-section">
  <h2 id="metadata-details-header">Details</h2>
  <div id='metadata'>
  <p>
    This changelog contains<strong> #{commits.length}</strong> commits<br/>
    Number of referenced commits is <strong>#{taskReferenceCount}</strong> which is <strong>#{healthgauge.round(1)}%</strong> of all commits<br/>
    Which leaves out <strong>#{unspecifiedCommitCount}</strong> commits without proper commit messages<br/>
    #{footer_html}
  </p>
  </div>
</div>
      EOM
    end

    #Returns a value representing the percentage of referenced commits
    def health(commits)
      taskReferenceCount = commits.values.select { |e| !@filterpatterns.match(e).nil? }.length
      unspecifiedCommitCount = commits.length - taskReferenceCount
      return (taskReferenceCount / commits.length.to_f)
    end

    #Takes a series of unique 'identifiers'. Creates a hash map where the key of a particular case points to a series of identifiers pointing to other stuff
    #{
    # "9923333"=>
    #  {
    #    "85c76dd6a4037085a5b0f6c986e392d4386395cd"=>"Second commit for fixed case 9923333", 
    #    "56bf0ea3ef12676c3bf555b24224e50695e10440"=>"Second commit for case 9923333", 
    #    "72e16e848926d491704236571e6e06b89bd8eed3"=>"Fixed some case 9923333"
    #  }, 
    # "9923838"=>
    #  {
    #    "37aed128ab4185e19bd5b5cca45a4c3aa1e6b963"=>"Fixed some case 9923838"
    #  }
    #}
    def task_id_list(commits)
      grouped_by_task_id = Hash.new

      #TODO: How do i share this filterpatterns?
      regex_arr = []
      @settings[:none]['regex'].each do |rx|
        regex_arr.push( eval(rx) )
      end


      if @settings[:none].has_key?('delimiter')
        split_pattern = eval(@settings[:none]['delimiter']) 
      end 

      @filterpatterns = Regexp.union(regex_arr) 

      commits.each do |k,v| 
         #Return the resulting ID of the regex
        match = @filterpatterns.match(v)
        # we require regexp to return group id if matches
        if !match.nil?
          # we could match empty string...
          if !match[:id].nil? or !match[:id].empty? or !match[:id].gsub!(/\s+/, "").empty?  then
            res = match[:id]
             if not split_pattern.nil?
              res.split(split_pattern).each do |split_value|
                if(!grouped_by_task_id.has_key?(split_value))
                  grouped_by_task_id[split_value] = Hash.new
                end
                grouped_by_task_id[split_value][k] = v
              end
            else
              if(!grouped_by_task_id.has_key?(res))
                grouped_by_task_id[res] = Hash.new
              end
              grouped_by_task_id[res][k] = v
            end
          end
        end 
      end
      grouped_by_task_id
    end

    # Returns the netto list of SHA that have not references.
    # Expected input for "all_commits"
    # commit list contains all SHAs, eg.
        #        {"79cbf45541f76f00e3b1a96b028e66fe65e680b9"=>
        #          "Test for multiple\n\nIssue: 1,2\n",
        #         "32ca499e0a0d022c6449b6f9a5b436213541c6b1"=>"Test for empty\n",
        #         "c10036985d6c3ba8892bd532cb5e7a9bca3952ee"=>
        #          "Test for none reference\n\nIssue: none\n",
        #         "e62d01c7548117dbf529f2eb93c448b3d1d865a9"=>
        #          "Updated readme file again - third commit\n\nIssue: 1\n",
        #         "19ae102a802dce3a8f457891a2da18ff43d75815"=>
        #          "Revert \"Updated readme file\"\n\nThis reverts commit 694946d80d69703cf70e9d77d58fca879157c158.\nIssue: 1\n",
        #         "694946d80d69703cf70e9d77d58fca879157c158"=>
        #          "Updated readme file\n\nIssue: 3\n",
        #         "7047925b7a29fac846c2842973fe828c22ac3a51"=>"Initial commit - added README\n"}
        # The tasks list above contain all our SHAs that have task references, but in a bit different list:
        #        {"1"=>
        #          {"79cbf45541f76f00e3b1a96b028e66fe65e680b9"=>
        #            "Test for multiple\n\nIssue: 1,2\n",
        #           "e62d01c7548117dbf529f2eb93c448b3d1d865a9"=>
        #            "Updated readme file again - third commit\n\nIssue: 1\n",
        #           "19ae102a802dce3a8f457891a2da18ff43d75815"=>
        #            "Revert \"Updated readme file\"\n\nThis reverts commit 694946d80d69703cf70e9d77d58fca879157c158.\nIssue: 1\n"},
        #         "2"=>
        #          {"79cbf45541f76f00e3b1a96b028e66fe65e680b9"=>
        #            "Test for multiple\n\nIssue: 1,2\n"},
        #         "3"=>
        #          {"694946d80d69703cf70e9d77d58fca879157c158"=>
        #            "Updated readme file\n\nIssue: 3\n"}}
        #
        # We want to diff those list, to generate a list of SHA that did reference a task
    def get_shas_without_reference(all_commits, task_commits)
      sha_list = Array.new
      task_commits.values.each {|array_entry| sha_list << array_entry.keys }
      task_shas = sha_list.flatten

      commit_shas = all_commits.keys

      # netto list of commit without tasks references:
      unreferenced_shas = commit_shas - task_shas
      return unreferenced_shas
    end

    #This one has become a tad too specific for git. We get the short shar (first 8 digits). And
    #we trim the subject (That is the first element in the list when splitting the message on a newline.  
    def task_id_report(commits)
      tasks = task_id_list(commits)
      mdPath = @path.nil? == true ? "ids.md" : File.join(@path, "#{@settings[:general]['id_log_name']}.md")

      File.open(mdPath,'w:UTF-8') do |file| 
        file << "\#PAC id report\n\n"
        tasks.each do |k,v|
          file << "\#\##{k}\n"
          v.each do |k1,v2| 
            file << " - #{k1.slice(0..7)}: #{v2.split(/\n/).first}\n"
          end
          file << "\n"          
        end
        file << "\#\#Unspecified\n"
        

        unreference_shas = get_shas_without_reference(commits, tasks)
        unreference_shas.each do |sha|
          file << " - #{ sha.slice(0..7) }: #{ commits[sha].split(/\n/).first }\n"
        end

        file << "\n"
        file << statistics_html(commits)
      end

      unless @settings[:general]['changelog_formats'].nil?
        if @settings[:general]['changelog_formats'].include?("html")
          html = Kramdown::Document.new(File.read(mdPath)).to_html
        end
        if @settings[:general]['changelog_formats'].include?("pdf")
           output_file_path = mdPath.sub(/\.md$/, '.pdf')
           kit = PDFKit.new(html, :page_size => 'A4')
           kit.to_file(output_file_path)
        end
      end
    end
  end

  #Very simple task system.
  class NoneTaskSystem < TaskSystem
    attr_accessor :filterpatterns
    def intialize(settings) 
      super(settings) 
    end
  end

  class FogBugzTaskSystem < TaskSystem

    attr_accessor :instance, :fields

    def initialize(settings)
      super(settings)
      @instance = Fogbugz::Interface.new( :email => @settings[:fogbugz]["fogbugz_usr"], :password => @settings[:fogbugz]["fogbugz_pwd"], :uri => @settings[:fogbugz]["fogbugz_url"] )
      
      regex_arr = []
      @settings[:fogbugz]['regex'].each do |rx|
        regex_arr.push( eval(rx) )
      end
      @filterpatterns = Regexp.union(regex_arr)      
      @instance.authenticate
    end

    #Filter commits method. This is wehre we apply our regexes.
    def filter(cmits)
      caseNumbers = []
      titleSearchStrings = []

      cmits.each do |cMsg|
        
        res = @filterpatterns.match(cMsg)
        if !res.nil?
          if res[:id].is_number?
            caseNumbers << res[:id]
          else
            titleSearchStrings << res[:id]
          end
        end
      end
      
      cases_xml_array = []
      #If the results of the parsing yields something that is not a case number (IE 1200, 2332 etc). Each item has to seperated with OR.
      if titleSearchStrings.length > 0
        xmloutput_title = @instance.command(:search, :q => "title:"+titleSearchStrings.uniq.join(" OR "), :cols => @settings[:fogbugz]["fogbugz_fields"])
        cases_xml_array.push(xmloutput_title)
      end
      
      
      if caseNumbers.length > 0
        caseNumbers.each do |c|
          case_xml = @instance.command(:search, :q => c, :cols => @settings[:fogbugz]["fogbugz_fields"])
          cases_xml_array.push(case_xml)
        end
      end
      
      hashes = []
      unless cases_xml_array.empty?
        cases_xml_array.each do |x|
          if x["cases"]["count"].to_i > 0
             if x["cases"]["case"].class == Hash
		h = Hash.new
                x["cases"]["case"].each do |aKey, aValue|
                  h[aKey.to_sym] = aValue
                end
                hashes.push(h)
             else
               x["cases"]["case"].each do |value|
		 h = Hash.new
                 value.each do |k,v|
                   h[k.to_sym] = v                    
                 end
                 hashes.push(h)
               end
             end  
          end
        end
      end
      descending = -1
      hashes.uniq { |x| x[:ixBug] }.sort_by { |hv| hv[:ixBug].to_i * descending }
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
              unless task[:sReleaseNotes].nil?
                file << "\n"
                file << "**Release note:** #{task[:sReleaseNotes]}\n"
                file << "\n"
                file << "**Status:** #{task[:sStatus]}\n"
                file << "\n"
                file << "**Title:** #{task[:sTitle]}\n"
                file << "\n"
              end
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
            unless task[:sStatus].downcase.include? "duplicate"
              unless task[:sReleaseNotes].nil? 
                file << "<div class ='change-item'>"
                file << ""

                title_case = "<div class='change-title' > <div class='issue-type-#{task[:sCategory]} issue-type'>&nbsp;</div>#{task[:sTitle]} (#{case_link_html(task[:ixBug])})</div>"
                title_case_escaped = html_escape_non_ascii(title_case)
                file << title_case_escaped

              
                rel_note = "<div class='release-note'>#{task[:sReleaseNotes]} </div>"
                rel_note_escaped = html_escape_non_ascii(rel_note)
                file << rel_note_escaped
              
                file << '</div>'
              end
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
