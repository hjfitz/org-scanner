// json in golang is a joke. This util should format the file better for the go script

const filename = process.argv[2]

try {
	const file = require(filename)
	const files = (file.files || []).map(fname => Object.keys(fname)[0])
	const security_issues = Object.keys(file.sec_issues || {}).map(key => file.sec_issues[key]).flat()
	const header_issues = Object.keys(file.missing_sec_header).map(key =>  file.missing_sec_header[key]).flat()
	console.log(JSON.stringify({
		files,
		security_issues,
		header_issues,
	}))
} catch (err) {
	console.log({})
}
