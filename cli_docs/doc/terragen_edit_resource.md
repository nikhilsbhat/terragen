## terragen edit resource

Command to edit already created scaffolds generated scaffolds of resource

### Synopsis

This will help user to edit scaffolds of resource that was already generated.
               Not all aspects of resource can be edited, it is very limited

```
terragen edit resource [args] [flags]
```

### Options

```
  -h, --help   help for resource
```

### Options inherited from parent commands

```
      --dry-run                dry-run the process of provider scaffold creation
  -f, --force                  enable this to forcefully create resource/datasource/importers (this might tamper the scaffold)
  -l, --log-level string       log level for terragen, log levels supported by [https://github.com/sirupsen/logrus] will work (default "info")
  -p, --path string            path where the templates has to be generated (default ".")
      --skip-provider-update   when enabled, updating provider.go with newly created datasource/resource would be skipped
      --use-plugin-framework   enable this to generate scaffolds with terraform-plugin-framework(https://github.com/hashicorp/terraform-plugin-framework)
```

### SEE ALSO

* [terragen edit](terragen_edit.md)	 - Command to edit the scaffold created for a provider

###### Auto generated by spf13/cobra on 23-Sep-2023
