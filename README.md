# Org Scanner

Use this to scan an entire GitHub organisation for:
1. Access keys
2. Unsafe dependencies
3. Static source code analysis (JavaScript only)

# Usage

1. Clone: `git clone git@github.com:hjfitz/org-scanner.git`
2. Install Node dependencies: `npm install` (or `yarn`)
3. Set up environment: `echo GITHUB_ACCESS_TOKEN=$MYGHACCESSTOKEN>.env`
4. Run: `node list-and-clone`

## One Repo
If you want to scan one repo, you can forgo a lot of the setup. Simply use **scan.sh**:

```bash
~ $ ./scan.sh $REPO_URL $ACCESS_TOKEN $REPO_NAME
```