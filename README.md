# cdkpw - A Profile-Aware Wrapper for AWS CDK

CDK

[![CI](https://github.com/phelian/cdkpw/actions/workflows/test.yml/badge.svg)](https://github.com/phelian/cdkpw/actions/workflows/test.yml)

**cdkpw** wraps the AWS CDK CLI and injects the right aws sso `--profile` based on your stack name. No more guessing which profile goes with which environment.

Example usage
`cdk deploy "Prod*"` - Would append --prod_admin based on the example configuration.

---

## ✨ Features

- 🧠 Auto-matches profile by stack name
- 🔒 Supports multiple profiles (dev, prod, secure, etc.)
- ⚙️ Configurable via YAML
- 💨 Drop-in replacement for `cdk`

## 🛠️ Config

Default: ~/.cdk/.cdkpw.yml
Or set CDKPW_CONFIG=/path/to/.cdkpw.yml

Example:

```yaml
profiles:
  - match: Prod
    profile: prod_admin
  - match: Dev
    profile: dev_admin
  - match: Secure
    profile: secure_admin
cdkLocation: ${CDK_BIN}
verbose: 0|1|2
```

cdkLocation defaults to `cdk`m accepts string or envvars  
verbose default to 0 (silent)

Verbose levels:

- 0 (Silent)
- 1 (Info)

## How to use

`alias cdk='cdkpw'` if alias is possible, if used in another nix flake i had to place a tiny cdk executable in front of path:

```bash
#!/usr/bin/env sh
exec cdkpw "$@"
```

## 📄 License

MIT — do whatever. Just don’t sue the author.

## ✍️ Author

Made by Phelian
PRs and stars welcome 🌟
