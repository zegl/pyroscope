require "pyroscope"

STDOUT.sync = true

Pyroscope.configure do |config|
  config.app_name = "broken.ruby.app"
  config.server_address = "http://pyroscope:4040/"
  config.tags = {
    "forked" => "false",
  }
end

def work(label, seconds)
  i = 0
  st = Time.now
  last_tag_update = -1
  while st + seconds > Time.now
    second = ((Time.now - st)).floor
    if last_tag_update != second
      puts "updating tags for #{label}, second #{second}"
      # I call it a 500 times because it doesn't hang consistently enough
      500.times do
        Pyroscope.tag({ "second" => second.to_s })
      end
      last_tag_update = second
    end
    i += 1
  end
end

def forked_process
  puts "forked_process work start"
  work("forked_process", 15)
  puts "forked_process work end"
end

def old_process
  puts "old_process work start"
  work("old_process", 15)
  puts "old_process work end"
end

# print multiline string
puts <<-EOS
this program should:
* create a fork of itself
* old process will run function `bar` for 15 seconds
* forked process will change its tags to 'forked => true'
* forked process will run function `foo` for 15 seconds
* both will add tags { "second" => X } for every second X
* results in pyroscope should show 15 seconds for old process and 15 seconds for new process
* there should be no double-counting
---

currently the there's a couple of problems:
* forked process eventually hangs when you call Pyroscope.tag enough times — big issue
* the data from forked process doesn't show up in pyroscope — this is a smaller issue (but we should still fix it).
  It's okay if users have to do extra set up work to make it work (e.g add Pyroscope.configure in `fork` method block)

EOS

puts "pre fork"
pid = fork do
  Pyroscope.tag({ "forked" => "true" })
  forked_process
end
puts "post fork"

old_process

Process.waitpid2(pid)
