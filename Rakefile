desc "build"
task :build do
  sh "rm -rf release/*"
  sh 'gox -os="darwin linux" -arch="386 amd64" -output "release/stns_{{.OS}}_{{.Arch}}/{{.Dir}}"'
end

desc "release"
task :release => :build do
  sh "ls release/ | xargs -o -I% zip -j release/%.zip release/%/stns-passwd"
  sh "ls release/ | grep -v zip | xargs -o -I% rm -rf release/%"
  v = `cat version.go | grep -i version | awk -F\\" '{ print $2}'`
  sh "ghr -u STNS -r stns-passwd #{v.chomp!} release/"
end
