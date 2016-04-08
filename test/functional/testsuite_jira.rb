module PAC__TestSuites
    #TestCase classes that contain the methods I want to use in the test suite
    require_relative 'testCases/jira_functional.rb'
    include PAC__TestCases_Jira
end