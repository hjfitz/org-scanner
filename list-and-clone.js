require('dotenv').config()
const {exec} = require('child_process')

const {Octokit} = require('@octokit/rest')

const accessToken = process.env.GITHUB_ACCESS_TOKEN

const client = new Octokit({
    auth: accessToken,
    userAgent: 'Security Scanner v0.0.0',
    timeZone: 'Europe/London',
    baseUrl: 'https://api.github.com',
})

const sh = cmd => new Promise((res, rej) => exec(cmd, (err, stdout) => err ? rej(err) : res(stdout)))

async function fetchRepos() { 
    let doPull = true
    const allRepos = []
    let page = 0;
    while (doPull) {
        console.log(`pulling from page ${page}`)
        const {data} = await client.repos
        .listForOrg({
            org: "elilillyco",
            per_page: 100,
            page,
        })
        allRepos.push(...data)
        doPull = !!data.length
        page += 1
    }
    return allRepos
}

async function scanRepos() {
    let repos;
    if (process.env.DEBUG_LOCAL === 'true') {
        repos = require('./static-data.json')
    } else {
        repos = await fetchRepos()
    }

    const elancoRepos = repos.filter(repo => repo.name.toLowerCase().includes('elanco'))
    
    
    for (const repo of elancoRepos) {     
        console.log(`scanning ${repo.name}`)
        try {  
            await sh(`./scan.sh ${repo.clone_url} ${process.env.GITHUB_ACCESS_TOKEN} ${repo.name}`)
        } catch (err) {
            console.error('error scanning:', err)
        }
    }
}

scanRepos()