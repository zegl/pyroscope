# README

This is a simple demonstration of using pyroscope gem with rubyonrails app.

In order to run this demo you may just run it with 

```
docker-compose up -d
```

Then you may point your browser to localhost:4040 to check generated flamegraphs.

Key changes you may find at `./config/application.rb`

```
   Pyroscope.configure do |config|
      config.app_name = "ride-sharing-app"
      config.server_address = "http://pyroscope:4040"
      config.tags = {
        "region": ENV["REGION"] || "us-east-1",
      }
```

As you may see here we define app name, server name and configure static tags for pyroscope agent instance.

At `./app/helpers/application_helper.rb` you may find dynamic tagging for each request.  

```
def find_nearest_vehicle(n, vehicle)
    Pyroscope.tag_wrapper({ "vehicle" => vehicle }) do
      i = 0
      start_time = Time.new
      while Time.new - start_time < n do
        i += 1
      end

      check_driver_availability(n) if vehicle == "car"
    end
  end
```
