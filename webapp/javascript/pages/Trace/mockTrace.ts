export const mockData = {
  data: {
    services: [
      {
        name: 'load-generator',
        numberOfSpans: 2,
      },
      {
        name: 'ride-sharing-app-java',
        numberOfSpans: 3,
      },
    ],
    spans: [
      {
        traceID: '1158f70ebb2bc1fb5084852761cecaae',
        spanID: '59eba9585b71d714',
        flags: 1,
        operationName: 'OrderVehicle',
        references: [],
        startTime: 1657875042429600,
        duration: 1515149,
        tags: [
          {
            key: 'internal.span.format',
            type: 'string',
            value: 'jaeger',
          },
          {
            key: 'otel.library.name',
            type: 'string',
            value: 'go.opentelemetry.io/otel/sdk/tracer',
          },
          {
            key: 'pyroscope.profile.baseline.url',
            type: 'string',
            value:
              'http://localhost:4040/comparison?from=1657871442&leftFrom=1657871442&leftQuery=load-generator.cpu%7Bspan_name%3D%22OrderVehicle%22%7D&leftUntil=1657875043&query=load-generator.cpu%7Bspan_name%3D%22OrderVehicle%22%7D&rightFrom=1657871442&rightQuery=load-generator.cpu%7Bprofile_id%3D%2259eba9585b71d714%22%7D&rightUntil=1657875043&until=1657875043',
          },
          {
            key: 'pyroscope.profile.diff.url',
            type: 'string',
            value:
              'http://localhost:4040/comparison-diff?from=1657871442&leftFrom=1657871442&leftQuery=load-generator.cpu%7Bspan_name%3D%22OrderVehicle%22%7D&leftUntil=1657875043&query=load-generator.cpu%7Bspan_name%3D%22OrderVehicle%22%7D&rightFrom=1657871442&rightQuery=load-generator.cpu%7Bprofile_id%3D%2259eba9585b71d714%22%7D&rightUntil=1657875043&until=1657875043',
          },
          {
            key: 'pyroscope.profile.id',
            type: 'string',
            value: '59eba9585b71d714',
          },
          {
            key: 'pyroscope.profile.url',
            type: 'string',
            value:
              'http://localhost:4040/?from=1657875042429620300&query=load-generator.cpu%7Bprofile_id%3D%2259eba9585b71d714%22%7D&until=1657875043944684700',
          },
          {
            key: 'vehicle',
            type: 'string',
            value: 'car',
          },
        ],
        logs: [],
        processID: 'p1',
        warnings: [],
        process: {
          serviceName: 'load-generator',
          tags: [],
        },
        relativeStartTime: 0,
        depth: 0,
        hasChildren: true,
      },
      {
        traceID: '1158f70ebb2bc1fb5084852761cecaae',
        spanID: '2fcaa4931dae97f1',
        flags: 1,
        operationName: 'HTTP GET',
        references: [
          {
            refType: 'CHILD_OF',
            traceID: '1158f70ebb2bc1fb5084852761cecaae',
            spanID: '59eba9585b71d714',
            span: {
              traceID: '1158f70ebb2bc1fb5084852761cecaae',
              spanID: '59eba9585b71d714',
              flags: 1,
              operationName: 'OrderVehicle',
              references: [],
              startTime: 1657875042429600,
              duration: 1515149,
              tags: [
                {
                  key: 'internal.span.format',
                  type: 'string',
                  value: 'jaeger',
                },
                {
                  key: 'otel.library.name',
                  type: 'string',
                  value: 'go.opentelemetry.io/otel/sdk/tracer',
                },
                {
                  key: 'pyroscope.profile.baseline.url',
                  type: 'string',
                  value:
                    'http://localhost:4040/comparison?from=1657871442&leftFrom=1657871442&leftQuery=load-generator.cpu%7Bspan_name%3D%22OrderVehicle%22%7D&leftUntil=1657875043&query=load-generator.cpu%7Bspan_name%3D%22OrderVehicle%22%7D&rightFrom=1657871442&rightQuery=load-generator.cpu%7Bprofile_id%3D%2259eba9585b71d714%22%7D&rightUntil=1657875043&until=1657875043',
                },
                {
                  key: 'pyroscope.profile.diff.url',
                  type: 'string',
                  value:
                    'http://localhost:4040/comparison-diff?from=1657871442&leftFrom=1657871442&leftQuery=load-generator.cpu%7Bspan_name%3D%22OrderVehicle%22%7D&leftUntil=1657875043&query=load-generator.cpu%7Bspan_name%3D%22OrderVehicle%22%7D&rightFrom=1657871442&rightQuery=load-generator.cpu%7Bprofile_id%3D%2259eba9585b71d714%22%7D&rightUntil=1657875043&until=1657875043',
                },
                {
                  key: 'pyroscope.profile.id',
                  type: 'string',
                  value: '59eba9585b71d714',
                },
                {
                  key: 'pyroscope.profile.url',
                  type: 'string',
                  value:
                    'http://localhost:4040/?from=1657875042429620300&query=load-generator.cpu%7Bprofile_id%3D%2259eba9585b71d714%22%7D&until=1657875043944684700',
                },
                {
                  key: 'vehicle',
                  type: 'string',
                  value: 'car',
                },
              ],
              logs: [],
              processID: 'p1',
              warnings: [],
              process: {
                serviceName: 'load-generator',
                tags: [],
              },
              relativeStartTime: 0,
              depth: 0,
              hasChildren: true,
            },
          },
        ],
        startTime: 1657875042735924,
        duration: 1208495,
        tags: [
          {
            key: 'http.flavor',
            type: 'string',
            value: '1.1',
          },
          {
            key: 'http.host',
            type: 'string',
            value: 'ap-south-java:5000',
          },
          {
            key: 'http.method',
            type: 'string',
            value: 'GET',
          },
          {
            key: 'http.scheme',
            type: 'string',
            value: 'http',
          },
          {
            key: 'http.status_code',
            type: 'int64',
            value: 200,
          },
          {
            key: 'http.url',
            type: 'string',
            value: 'http://ap-south-java:5000/car',
          },
          {
            key: 'internal.span.format',
            type: 'string',
            value: 'jaeger',
          },
          {
            key: 'otel.library.name',
            type: 'string',
            value:
              'go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp',
          },
          {
            key: 'otel.library.version',
            type: 'string',
            value: 'semver:0.27.0',
          },
          {
            key: 'span.kind',
            type: 'string',
            value: 'client',
          },
        ],
        logs: [],
        processID: 'p1',
        warnings: [],
        process: {
          serviceName: 'load-generator',
          tags: [],
        },
        relativeStartTime: 306324,
        depth: 1,
        hasChildren: true,
      },
      {
        traceID: '1158f70ebb2bc1fb5084852761cecaae',
        spanID: 'e4d3228f084958ad',
        operationName: 'orderCar',
        references: [
          {
            refType: 'CHILD_OF',
            traceID: '1158f70ebb2bc1fb5084852761cecaae',
            spanID: '2fcaa4931dae97f1',
            span: {
              traceID: '1158f70ebb2bc1fb5084852761cecaae',
              spanID: '2fcaa4931dae97f1',
              flags: 1,
              operationName: 'HTTP GET',
              references: [
                {
                  refType: 'CHILD_OF',
                  traceID: '1158f70ebb2bc1fb5084852761cecaae',
                  spanID: '59eba9585b71d714',
                  span: {
                    traceID: '1158f70ebb2bc1fb5084852761cecaae',
                    spanID: '59eba9585b71d714',
                    flags: 1,
                    operationName: 'OrderVehicle',
                    references: [],
                    startTime: 1657875042429600,
                    duration: 1515149,
                    tags: [
                      {
                        key: 'internal.span.format',
                        type: 'string',
                        value: 'jaeger',
                      },
                      {
                        key: 'otel.library.name',
                        type: 'string',
                        value: 'go.opentelemetry.io/otel/sdk/tracer',
                      },
                      {
                        key: 'pyroscope.profile.baseline.url',
                        type: 'string',
                        value:
                          'http://localhost:4040/comparison?from=1657871442&leftFrom=1657871442&leftQuery=load-generator.cpu%7Bspan_name%3D%22OrderVehicle%22%7D&leftUntil=1657875043&query=load-generator.cpu%7Bspan_name%3D%22OrderVehicle%22%7D&rightFrom=1657871442&rightQuery=load-generator.cpu%7Bprofile_id%3D%2259eba9585b71d714%22%7D&rightUntil=1657875043&until=1657875043',
                      },
                      {
                        key: 'pyroscope.profile.diff.url',
                        type: 'string',
                        value:
                          'http://localhost:4040/comparison-diff?from=1657871442&leftFrom=1657871442&leftQuery=load-generator.cpu%7Bspan_name%3D%22OrderVehicle%22%7D&leftUntil=1657875043&query=load-generator.cpu%7Bspan_name%3D%22OrderVehicle%22%7D&rightFrom=1657871442&rightQuery=load-generator.cpu%7Bprofile_id%3D%2259eba9585b71d714%22%7D&rightUntil=1657875043&until=1657875043',
                      },
                      {
                        key: 'pyroscope.profile.id',
                        type: 'string',
                        value: '59eba9585b71d714',
                      },
                      {
                        key: 'pyroscope.profile.url',
                        type: 'string',
                        value:
                          'http://localhost:4040/?from=1657875042429620300&query=load-generator.cpu%7Bprofile_id%3D%2259eba9585b71d714%22%7D&until=1657875043944684700',
                      },
                      {
                        key: 'vehicle',
                        type: 'string',
                        value: 'car',
                      },
                    ],
                    logs: [],
                    processID: 'p1',
                    warnings: [],
                    process: {
                      serviceName: 'load-generator',
                      tags: [],
                    },
                    relativeStartTime: 0,
                    depth: 0,
                    hasChildren: true,
                  },
                },
              ],
              startTime: 1657875042735924,
              duration: 1208495,
              tags: [
                {
                  key: 'http.flavor',
                  type: 'string',
                  value: '1.1',
                },
                {
                  key: 'http.host',
                  type: 'string',
                  value: 'ap-south-java:5000',
                },
                {
                  key: 'http.method',
                  type: 'string',
                  value: 'GET',
                },
                {
                  key: 'http.scheme',
                  type: 'string',
                  value: 'http',
                },
                {
                  key: 'http.status_code',
                  type: 'int64',
                  value: 200,
                },
                {
                  key: 'http.url',
                  type: 'string',
                  value: 'http://ap-south-java:5000/car',
                },
                {
                  key: 'internal.span.format',
                  type: 'string',
                  value: 'jaeger',
                },
                {
                  key: 'otel.library.name',
                  type: 'string',
                  value:
                    'go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp',
                },
                {
                  key: 'otel.library.version',
                  type: 'string',
                  value: 'semver:0.27.0',
                },
                {
                  key: 'span.kind',
                  type: 'string',
                  value: 'client',
                },
              ],
              logs: [],
              processID: 'p1',
              warnings: [],
              process: {
                serviceName: 'load-generator',
                tags: [],
              },
              relativeStartTime: 306324,
              depth: 1,
              hasChildren: true,
            },
          },
        ],
        startTime: 1657875042740200,
        duration: 1203011,
        tags: [
          {
            key: 'internal.span.format',
            type: 'string',
            value: 'proto',
          },
          {
            key: 'otel.library.name',
            type: 'string',
            value: 'ride-sharing-app-java',
          },
          {
            key: 'otel.scope.name',
            type: 'string',
            value: 'ride-sharing-app-java',
          },
          {
            key: 'pyroscope.profile.baseline.url',
            type: 'string',
            value:
              'http://localhost:4040/comparison?query=ride-sharing-app-java.itimer%7Bspan_name%3D%22orderCar%22%7D&from=1657871442741&until=1657875043943&leftQuery=ride-sharing-app-java.itimer%7Bspan_name%3D%22orderCar%22%7D&leftFrom=1657871442741&leftUntil=1657875043943&rightQuery=ride-sharing-app-java.itimer%7Bprofile_id%3D%22e4d3228f084958ad%22%7D&rightFrom=1657871442741&rightUntil=1657875043943',
          },
          {
            key: 'pyroscope.profile.diff.url',
            type: 'string',
            value:
              'http://localhost:4040/comparison-diff?query=ride-sharing-app-java.itimer%7Bspan_name%3D%22orderCar%22%7D&from=1657871442741&until=1657875043943&leftQuery=ride-sharing-app-java.itimer%7Bspan_name%3D%22orderCar%22%7D&leftFrom=1657871442741&leftUntil=1657875043943&rightQuery=ride-sharing-app-java.itimer%7Bprofile_id%3D%22e4d3228f084958ad%22%7D&rightFrom=1657871442741&rightUntil=1657875043943',
          },
          {
            key: 'pyroscope.profile.id',
            type: 'string',
            value: 'e4d3228f084958ad',
          },
          {
            key: 'pyroscope.profile.url',
            type: 'string',
            value:
              'http://localhost:4040?query=ride-sharing-app-java.itimer%7Bprofile_id%3D%22e4d3228f084958ad%22%7D&from=1657875042741&until=1657875043943',
          },
        ],
        logs: [],
        processID: 'p2',
        warnings: [],
        process: {
          serviceName: 'ride-sharing-app-java',
          tags: [
            {
              key: 'hostname',
              type: 'string',
              value: '9830275e739d',
            },
            {
              key: 'ip',
              type: 'string',
              value: '172.18.0.3',
            },
            {
              key: 'jaeger.version',
              type: 'string',
              value: 'opentelemetry-java',
            },
            {
              key: 'service.name',
              type: 'string',
              value: 'ride-sharing-app-java',
            },
            {
              key: 'telemetry.sdk.language',
              type: 'string',
              value: 'java',
            },
            {
              key: 'telemetry.sdk.name',
              type: 'string',
              value: 'opentelemetry',
            },
            {
              key: 'telemetry.sdk.version',
              type: 'string',
              value: '1.14.0',
            },
          ],
        },
        relativeStartTime: 310600,
        depth: 2,
        hasChildren: true,
      },
      {
        traceID: '1158f70ebb2bc1fb5084852761cecaae',
        spanID: '4b0404776651d6cc',
        operationName: 'findNearestVehicle',
        references: [
          {
            refType: 'CHILD_OF',
            traceID: '1158f70ebb2bc1fb5084852761cecaae',
            spanID: 'e4d3228f084958ad',
            span: {
              traceID: '1158f70ebb2bc1fb5084852761cecaae',
              spanID: 'e4d3228f084958ad',
              operationName: 'orderCar',
              references: [
                {
                  refType: 'CHILD_OF',
                  traceID: '1158f70ebb2bc1fb5084852761cecaae',
                  spanID: '2fcaa4931dae97f1',
                  span: {
                    traceID: '1158f70ebb2bc1fb5084852761cecaae',
                    spanID: '2fcaa4931dae97f1',
                    flags: 1,
                    operationName: 'HTTP GET',
                    references: [
                      {
                        refType: 'CHILD_OF',
                        traceID: '1158f70ebb2bc1fb5084852761cecaae',
                        spanID: '59eba9585b71d714',
                        span: {
                          traceID: '1158f70ebb2bc1fb5084852761cecaae',
                          spanID: '59eba9585b71d714',
                          flags: 1,
                          operationName: 'OrderVehicle',
                          references: [],
                          startTime: 1657875042429600,
                          duration: 1515149,
                          tags: [
                            {
                              key: 'internal.span.format',
                              type: 'string',
                              value: 'jaeger',
                            },
                            {
                              key: 'otel.library.name',
                              type: 'string',
                              value: 'go.opentelemetry.io/otel/sdk/tracer',
                            },
                            {
                              key: 'pyroscope.profile.baseline.url',
                              type: 'string',
                              value:
                                'http://localhost:4040/comparison?from=1657871442&leftFrom=1657871442&leftQuery=load-generator.cpu%7Bspan_name%3D%22OrderVehicle%22%7D&leftUntil=1657875043&query=load-generator.cpu%7Bspan_name%3D%22OrderVehicle%22%7D&rightFrom=1657871442&rightQuery=load-generator.cpu%7Bprofile_id%3D%2259eba9585b71d714%22%7D&rightUntil=1657875043&until=1657875043',
                            },
                            {
                              key: 'pyroscope.profile.diff.url',
                              type: 'string',
                              value:
                                'http://localhost:4040/comparison-diff?from=1657871442&leftFrom=1657871442&leftQuery=load-generator.cpu%7Bspan_name%3D%22OrderVehicle%22%7D&leftUntil=1657875043&query=load-generator.cpu%7Bspan_name%3D%22OrderVehicle%22%7D&rightFrom=1657871442&rightQuery=load-generator.cpu%7Bprofile_id%3D%2259eba9585b71d714%22%7D&rightUntil=1657875043&until=1657875043',
                            },
                            {
                              key: 'pyroscope.profile.id',
                              type: 'string',
                              value: '59eba9585b71d714',
                            },
                            {
                              key: 'pyroscope.profile.url',
                              type: 'string',
                              value:
                                'http://localhost:4040/?from=1657875042429620300&query=load-generator.cpu%7Bprofile_id%3D%2259eba9585b71d714%22%7D&until=1657875043944684700',
                            },
                            {
                              key: 'vehicle',
                              type: 'string',
                              value: 'car',
                            },
                          ],
                          logs: [],
                          processID: 'p1',
                          warnings: [],
                          process: {
                            serviceName: 'load-generator',
                            tags: [],
                          },
                          relativeStartTime: 0,
                          depth: 0,
                          hasChildren: true,
                        },
                      },
                    ],
                    startTime: 1657875042735924,
                    duration: 1208495,
                    tags: [
                      {
                        key: 'http.flavor',
                        type: 'string',
                        value: '1.1',
                      },
                      {
                        key: 'http.host',
                        type: 'string',
                        value: 'ap-south-java:5000',
                      },
                      {
                        key: 'http.method',
                        type: 'string',
                        value: 'GET',
                      },
                      {
                        key: 'http.scheme',
                        type: 'string',
                        value: 'http',
                      },
                      {
                        key: 'http.status_code',
                        type: 'int64',
                        value: 200,
                      },
                      {
                        key: 'http.url',
                        type: 'string',
                        value: 'http://ap-south-java:5000/car',
                      },
                      {
                        key: 'internal.span.format',
                        type: 'string',
                        value: 'jaeger',
                      },
                      {
                        key: 'otel.library.name',
                        type: 'string',
                        value:
                          'go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp',
                      },
                      {
                        key: 'otel.library.version',
                        type: 'string',
                        value: 'semver:0.27.0',
                      },
                      {
                        key: 'span.kind',
                        type: 'string',
                        value: 'client',
                      },
                    ],
                    logs: [],
                    processID: 'p1',
                    warnings: [],
                    process: {
                      serviceName: 'load-generator',
                      tags: [],
                    },
                    relativeStartTime: 306324,
                    depth: 1,
                    hasChildren: true,
                  },
                },
              ],
              startTime: 1657875042740200,
              duration: 1203011,
              tags: [
                {
                  key: 'internal.span.format',
                  type: 'string',
                  value: 'proto',
                },
                {
                  key: 'otel.library.name',
                  type: 'string',
                  value: 'ride-sharing-app-java',
                },
                {
                  key: 'otel.scope.name',
                  type: 'string',
                  value: 'ride-sharing-app-java',
                },
                {
                  key: 'pyroscope.profile.baseline.url',
                  type: 'string',
                  value:
                    'http://localhost:4040/comparison?query=ride-sharing-app-java.itimer%7Bspan_name%3D%22orderCar%22%7D&from=1657871442741&until=1657875043943&leftQuery=ride-sharing-app-java.itimer%7Bspan_name%3D%22orderCar%22%7D&leftFrom=1657871442741&leftUntil=1657875043943&rightQuery=ride-sharing-app-java.itimer%7Bprofile_id%3D%22e4d3228f084958ad%22%7D&rightFrom=1657871442741&rightUntil=1657875043943',
                },
                {
                  key: 'pyroscope.profile.diff.url',
                  type: 'string',
                  value:
                    'http://localhost:4040/comparison-diff?query=ride-sharing-app-java.itimer%7Bspan_name%3D%22orderCar%22%7D&from=1657871442741&until=1657875043943&leftQuery=ride-sharing-app-java.itimer%7Bspan_name%3D%22orderCar%22%7D&leftFrom=1657871442741&leftUntil=1657875043943&rightQuery=ride-sharing-app-java.itimer%7Bprofile_id%3D%22e4d3228f084958ad%22%7D&rightFrom=1657871442741&rightUntil=1657875043943',
                },
                {
                  key: 'pyroscope.profile.id',
                  type: 'string',
                  value: 'e4d3228f084958ad',
                },
                {
                  key: 'pyroscope.profile.url',
                  type: 'string',
                  value:
                    'http://localhost:4040?query=ride-sharing-app-java.itimer%7Bprofile_id%3D%22e4d3228f084958ad%22%7D&from=1657875042741&until=1657875043943',
                },
              ],
              logs: [],
              processID: 'p2',
              warnings: [],
              process: {
                serviceName: 'ride-sharing-app-java',
                tags: [
                  {
                    key: 'hostname',
                    type: 'string',
                    value: '9830275e739d',
                  },
                  {
                    key: 'ip',
                    type: 'string',
                    value: '172.18.0.3',
                  },
                  {
                    key: 'jaeger.version',
                    type: 'string',
                    value: 'opentelemetry-java',
                  },
                  {
                    key: 'service.name',
                    type: 'string',
                    value: 'ride-sharing-app-java',
                  },
                  {
                    key: 'telemetry.sdk.language',
                    type: 'string',
                    value: 'java',
                  },
                  {
                    key: 'telemetry.sdk.name',
                    type: 'string',
                    value: 'opentelemetry',
                  },
                  {
                    key: 'telemetry.sdk.version',
                    type: 'string',
                    value: '1.14.0',
                  },
                ],
              },
              relativeStartTime: 310600,
              depth: 2,
              hasChildren: true,
            },
          },
        ],
        startTime: 1657875042741219,
        duration: 1201679,
        tags: [
          {
            key: 'internal.span.format',
            type: 'string',
            value: 'proto',
          },
          {
            key: 'otel.library.name',
            type: 'string',
            value: 'ride-sharing-app-java',
          },
          {
            key: 'otel.scope.name',
            type: 'string',
            value: 'ride-sharing-app-java',
          },
        ],
        logs: [],
        processID: 'p2',
        warnings: [],
        process: {
          serviceName: 'ride-sharing-app-java',
          tags: [
            {
              key: 'hostname',
              type: 'string',
              value: '9830275e739d',
            },
            {
              key: 'ip',
              type: 'string',
              value: '172.18.0.3',
            },
            {
              key: 'jaeger.version',
              type: 'string',
              value: 'opentelemetry-java',
            },
            {
              key: 'service.name',
              type: 'string',
              value: 'ride-sharing-app-java',
            },
            {
              key: 'telemetry.sdk.language',
              type: 'string',
              value: 'java',
            },
            {
              key: 'telemetry.sdk.name',
              type: 'string',
              value: 'opentelemetry',
            },
            {
              key: 'telemetry.sdk.version',
              type: 'string',
              value: '1.14.0',
            },
          ],
        },
        relativeStartTime: 311619,
        depth: 3,
        hasChildren: true,
      },
      {
        traceID: '1158f70ebb2bc1fb5084852761cecaae',
        spanID: 'a03a2313125fd5c3',
        operationName: 'checkDriverAvailability',
        references: [
          {
            refType: 'CHILD_OF',
            traceID: '1158f70ebb2bc1fb5084852761cecaae',
            spanID: '4b0404776651d6cc',
            span: {
              traceID: '1158f70ebb2bc1fb5084852761cecaae',
              spanID: '4b0404776651d6cc',
              operationName: 'findNearestVehicle',
              references: [
                {
                  refType: 'CHILD_OF',
                  traceID: '1158f70ebb2bc1fb5084852761cecaae',
                  spanID: 'e4d3228f084958ad',
                  span: {
                    traceID: '1158f70ebb2bc1fb5084852761cecaae',
                    spanID: 'e4d3228f084958ad',
                    operationName: 'orderCar',
                    references: [
                      {
                        refType: 'CHILD_OF',
                        traceID: '1158f70ebb2bc1fb5084852761cecaae',
                        spanID: '2fcaa4931dae97f1',
                        span: {
                          traceID: '1158f70ebb2bc1fb5084852761cecaae',
                          spanID: '2fcaa4931dae97f1',
                          flags: 1,
                          operationName: 'HTTP GET',
                          references: [
                            {
                              refType: 'CHILD_OF',
                              traceID: '1158f70ebb2bc1fb5084852761cecaae',
                              spanID: '59eba9585b71d714',
                              span: {
                                traceID: '1158f70ebb2bc1fb5084852761cecaae',
                                spanID: '59eba9585b71d714',
                                flags: 1,
                                operationName: 'OrderVehicle',
                                references: [],
                                startTime: 1657875042429600,
                                duration: 1515149,
                                tags: [
                                  {
                                    key: 'internal.span.format',
                                    type: 'string',
                                    value: 'jaeger',
                                  },
                                  {
                                    key: 'otel.library.name',
                                    type: 'string',
                                    value:
                                      'go.opentelemetry.io/otel/sdk/tracer',
                                  },
                                  {
                                    key: 'pyroscope.profile.baseline.url',
                                    type: 'string',
                                    value:
                                      'http://localhost:4040/comparison?from=1657871442&leftFrom=1657871442&leftQuery=load-generator.cpu%7Bspan_name%3D%22OrderVehicle%22%7D&leftUntil=1657875043&query=load-generator.cpu%7Bspan_name%3D%22OrderVehicle%22%7D&rightFrom=1657871442&rightQuery=load-generator.cpu%7Bprofile_id%3D%2259eba9585b71d714%22%7D&rightUntil=1657875043&until=1657875043',
                                  },
                                  {
                                    key: 'pyroscope.profile.diff.url',
                                    type: 'string',
                                    value:
                                      'http://localhost:4040/comparison-diff?from=1657871442&leftFrom=1657871442&leftQuery=load-generator.cpu%7Bspan_name%3D%22OrderVehicle%22%7D&leftUntil=1657875043&query=load-generator.cpu%7Bspan_name%3D%22OrderVehicle%22%7D&rightFrom=1657871442&rightQuery=load-generator.cpu%7Bprofile_id%3D%2259eba9585b71d714%22%7D&rightUntil=1657875043&until=1657875043',
                                  },
                                  {
                                    key: 'pyroscope.profile.id',
                                    type: 'string',
                                    value: '59eba9585b71d714',
                                  },
                                  {
                                    key: 'pyroscope.profile.url',
                                    type: 'string',
                                    value:
                                      'http://localhost:4040/?from=1657875042429620300&query=load-generator.cpu%7Bprofile_id%3D%2259eba9585b71d714%22%7D&until=1657875043944684700',
                                  },
                                  {
                                    key: 'vehicle',
                                    type: 'string',
                                    value: 'car',
                                  },
                                ],
                                logs: [],
                                processID: 'p1',
                                warnings: [],
                                process: {
                                  serviceName: 'load-generator',
                                  tags: [],
                                },
                                relativeStartTime: 0,
                                depth: 0,
                                hasChildren: true,
                              },
                            },
                          ],
                          startTime: 1657875042735924,
                          duration: 1208495,
                          tags: [
                            {
                              key: 'http.flavor',
                              type: 'string',
                              value: '1.1',
                            },
                            {
                              key: 'http.host',
                              type: 'string',
                              value: 'ap-south-java:5000',
                            },
                            {
                              key: 'http.method',
                              type: 'string',
                              value: 'GET',
                            },
                            {
                              key: 'http.scheme',
                              type: 'string',
                              value: 'http',
                            },
                            {
                              key: 'http.status_code',
                              type: 'int64',
                              value: 200,
                            },
                            {
                              key: 'http.url',
                              type: 'string',
                              value: 'http://ap-south-java:5000/car',
                            },
                            {
                              key: 'internal.span.format',
                              type: 'string',
                              value: 'jaeger',
                            },
                            {
                              key: 'otel.library.name',
                              type: 'string',
                              value:
                                'go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp',
                            },
                            {
                              key: 'otel.library.version',
                              type: 'string',
                              value: 'semver:0.27.0',
                            },
                            {
                              key: 'span.kind',
                              type: 'string',
                              value: 'client',
                            },
                          ],
                          logs: [],
                          processID: 'p1',
                          warnings: [],
                          process: {
                            serviceName: 'load-generator',
                            tags: [],
                          },
                          relativeStartTime: 306324,
                          depth: 1,
                          hasChildren: true,
                        },
                      },
                    ],
                    startTime: 1657875042740200,
                    duration: 1203011,
                    tags: [
                      {
                        key: 'internal.span.format',
                        type: 'string',
                        value: 'proto',
                      },
                      {
                        key: 'otel.library.name',
                        type: 'string',
                        value: 'ride-sharing-app-java',
                      },
                      {
                        key: 'otel.scope.name',
                        type: 'string',
                        value: 'ride-sharing-app-java',
                      },
                      {
                        key: 'pyroscope.profile.baseline.url',
                        type: 'string',
                        value:
                          'http://localhost:4040/comparison?query=ride-sharing-app-java.itimer%7Bspan_name%3D%22orderCar%22%7D&from=1657871442741&until=1657875043943&leftQuery=ride-sharing-app-java.itimer%7Bspan_name%3D%22orderCar%22%7D&leftFrom=1657871442741&leftUntil=1657875043943&rightQuery=ride-sharing-app-java.itimer%7Bprofile_id%3D%22e4d3228f084958ad%22%7D&rightFrom=1657871442741&rightUntil=1657875043943',
                      },
                      {
                        key: 'pyroscope.profile.diff.url',
                        type: 'string',
                        value:
                          'http://localhost:4040/comparison-diff?query=ride-sharing-app-java.itimer%7Bspan_name%3D%22orderCar%22%7D&from=1657871442741&until=1657875043943&leftQuery=ride-sharing-app-java.itimer%7Bspan_name%3D%22orderCar%22%7D&leftFrom=1657871442741&leftUntil=1657875043943&rightQuery=ride-sharing-app-java.itimer%7Bprofile_id%3D%22e4d3228f084958ad%22%7D&rightFrom=1657871442741&rightUntil=1657875043943',
                      },
                      {
                        key: 'pyroscope.profile.id',
                        type: 'string',
                        value: 'e4d3228f084958ad',
                      },
                      {
                        key: 'pyroscope.profile.url',
                        type: 'string',
                        value:
                          'http://localhost:4040?query=ride-sharing-app-java.itimer%7Bprofile_id%3D%22e4d3228f084958ad%22%7D&from=1657875042741&until=1657875043943',
                      },
                    ],
                    logs: [],
                    processID: 'p2',
                    warnings: [],
                    process: {
                      serviceName: 'ride-sharing-app-java',
                      tags: [
                        {
                          key: 'hostname',
                          type: 'string',
                          value: '9830275e739d',
                        },
                        {
                          key: 'ip',
                          type: 'string',
                          value: '172.18.0.3',
                        },
                        {
                          key: 'jaeger.version',
                          type: 'string',
                          value: 'opentelemetry-java',
                        },
                        {
                          key: 'service.name',
                          type: 'string',
                          value: 'ride-sharing-app-java',
                        },
                        {
                          key: 'telemetry.sdk.language',
                          type: 'string',
                          value: 'java',
                        },
                        {
                          key: 'telemetry.sdk.name',
                          type: 'string',
                          value: 'opentelemetry',
                        },
                        {
                          key: 'telemetry.sdk.version',
                          type: 'string',
                          value: '1.14.0',
                        },
                      ],
                    },
                    relativeStartTime: 310600,
                    depth: 2,
                    hasChildren: true,
                  },
                },
              ],
              startTime: 1657875042741219,
              duration: 1201679,
              tags: [
                {
                  key: 'internal.span.format',
                  type: 'string',
                  value: 'proto',
                },
                {
                  key: 'otel.library.name',
                  type: 'string',
                  value: 'ride-sharing-app-java',
                },
                {
                  key: 'otel.scope.name',
                  type: 'string',
                  value: 'ride-sharing-app-java',
                },
              ],
              logs: [],
              processID: 'p2',
              warnings: [],
              process: {
                serviceName: 'ride-sharing-app-java',
                tags: [
                  {
                    key: 'hostname',
                    type: 'string',
                    value: '9830275e739d',
                  },
                  {
                    key: 'ip',
                    type: 'string',
                    value: '172.18.0.3',
                  },
                  {
                    key: 'jaeger.version',
                    type: 'string',
                    value: 'opentelemetry-java',
                  },
                  {
                    key: 'service.name',
                    type: 'string',
                    value: 'ride-sharing-app-java',
                  },
                  {
                    key: 'telemetry.sdk.language',
                    type: 'string',
                    value: 'java',
                  },
                  {
                    key: 'telemetry.sdk.name',
                    type: 'string',
                    value: 'opentelemetry',
                  },
                  {
                    key: 'telemetry.sdk.version',
                    type: 'string',
                    value: '1.14.0',
                  },
                ],
              },
              relativeStartTime: 311619,
              depth: 3,
              hasChildren: true,
            },
          },
        ],
        startTime: 1657875043342078,
        duration: 600785,
        tags: [
          {
            key: 'internal.span.format',
            type: 'string',
            value: 'proto',
          },
          {
            key: 'otel.library.name',
            type: 'string',
            value: 'ride-sharing-app-java',
          },
          {
            key: 'otel.scope.name',
            type: 'string',
            value: 'ride-sharing-app-java',
          },
        ],
        logs: [],
        processID: 'p2',
        warnings: [],
        process: {
          serviceName: 'ride-sharing-app-java',
          tags: [
            {
              key: 'hostname',
              type: 'string',
              value: '9830275e739d',
            },
            {
              key: 'ip',
              type: 'string',
              value: '172.18.0.3',
            },
            {
              key: 'jaeger.version',
              type: 'string',
              value: 'opentelemetry-java',
            },
            {
              key: 'service.name',
              type: 'string',
              value: 'ride-sharing-app-java',
            },
            {
              key: 'telemetry.sdk.language',
              type: 'string',
              value: 'java',
            },
            {
              key: 'telemetry.sdk.name',
              type: 'string',
              value: 'opentelemetry',
            },
            {
              key: 'telemetry.sdk.version',
              type: 'string',
              value: '1.14.0',
            },
          ],
        },
        relativeStartTime: 912478,
        depth: 4,
        hasChildren: false,
      },
    ],
    traceID: '1158f70ebb2bc1fb5084852761cecaae',
    traceName: 'load-generator: OrderVehicle',
    processes: {
      p1: {
        serviceName: 'load-generator',
        tags: [],
      },
      p2: {
        serviceName: 'ride-sharing-app-java',
        tags: [
          {
            key: 'hostname',
            type: 'string',
            value: '9830275e739d',
          },
          {
            key: 'ip',
            type: 'string',
            value: '172.18.0.3',
          },
          {
            key: 'jaeger.version',
            type: 'string',
            value: 'opentelemetry-java',
          },
          {
            key: 'service.name',
            type: 'string',
            value: 'ride-sharing-app-java',
          },
          {
            key: 'telemetry.sdk.language',
            type: 'string',
            value: 'java',
          },
          {
            key: 'telemetry.sdk.name',
            type: 'string',
            value: 'opentelemetry',
          },
          {
            key: 'telemetry.sdk.version',
            type: 'string',
            value: '1.14.0',
          },
        ],
      },
    },
    duration: 1515149,
    startTime: 1657875042429600,
    endTime: 1657875043944749,
  },
  id: '1158f70ebb2bc1fb5084852761cecaae',
  state: 'FETCH_DONE',
};
