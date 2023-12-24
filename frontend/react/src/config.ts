interface AppConfig {
  domain: string;
  apiBaseUrl: string;
  minioUrl: string;
}

const domain: string = 'staging.game.k8s.atpstealer.com'
const subdomain: string = 'api.'
const port: string = ''

// const domain: string = 'localhost'
// const subdomain: string = ''
// const port: string = ':8000'

const config: AppConfig = {
  domain,
  apiBaseUrl: 'http://' + subdomain + domain + port + '/v1',
  minioUrl: 'https://minio.game.k8s.atpstealer.com/game-staging/assets'
}

export default config