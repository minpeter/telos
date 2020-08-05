import path from 'path'
import fs from 'fs'
import yaml from 'yaml'
import deepMerge from 'deepmerge'
import { PartialDeep } from 'type-fest'
import { ACL } from '../util/restrict'

export type ProviderConfig = {
  name: string;
  options: unknown;
}

export type Sponsor = {
  name: string;
  icon: string;
  description: string;
  small?: boolean;
}

export type ServerConfig = {
  database: {
    sql: string | {
      host: string;
      port: number;
      user: string;
      password: string;
      database: string;
    }
    redis: string | {
      host: string;
      port: number;
      password: string;
      database: number;
    }
    migrate: string;
  }
  instanceType: string;
  tokenKey: string;
  origin: string;

  ctftime?: {
    clientId: string;
    clientSecret: string;
  }

  userMembers: boolean;
  sponsors: Sponsor[];
  homeContent: string;
  ctfName: string;
  meta: {
    description: string;
    imageUrl: string;
  }
  logoUrl?: string;
  globalSiteTag?: string;

  challengeProvider: ProviderConfig;
  uploadProvider: ProviderConfig;

  email?: {
    provider: ProviderConfig;
    from: string;
  }

  divisions: Record<string, string>;
  defaultDivision?: string;
  divisionACLs?: ACL[];

  startTime: number;
  endTime: number;

  leaderboard: {
    maxLimit: number;
    maxOffset: number;
    updateInterval: number;
    graphMaxTeams: number;
    graphSampleTime: number;
  }
  loginTimeout: number;
}

const fileConfigLoaders: Map<string, (file: string) => unknown> = new Map([
  ['json', file => JSON.parse(file)],
  ['yaml', file => yaml.parse(file)],
  ['yml', file => yaml.parse(file)]
])

const configPath = path.resolve(__dirname, '../../..', process.env.RCTF_CONF_PATH ?? 'conf.d')
const fileConfigs: PartialDeep<ServerConfig>[] = []
fs.readdirSync(configPath).sort().forEach((name) => {
  const matched = name.match(/\.(.+)$/)
  if (matched === null || !fileConfigLoaders.has(matched[1])) {
    return
  }
  const loader = fileConfigLoaders.get(matched[1])
  const config = loader(fs.readFileSync(path.join(configPath, name)).toString())
  fileConfigs.push(config)
})

const parseBoolEnv = (val?: string): boolean | null => {
  if (val == null) {
    return null
  }
  return ['true', 'yes', 'y', '1'].includes(val.toLowerCase().trim())
}

const envConfig: ServerConfig = {
  database: {
    sql: process.env.RCTF_DATABASE_URL ?? {
      host: process.env.RCTF_DATABASE_HOST,
      port: parseInt(process.env.RCTF_DATABASE_PORT) || undefined,
      user: process.env.RCTF_DATABASE_USERNAME,
      password: process.env.RCTF_DATABASE_PASSWORD,
      database: process.env.RCTF_DATABASE_DATABASE
    },
    redis: process.env.RCTF_REDIS_URL ?? {
      host: process.env.RCTF_REDIS_HOST,
      port: parseInt(process.env.RCTF_REDIS_PORT) || undefined,
      password: process.env.RCTF_REDIS_PASSWORD,
      database: parseInt(process.env.RCTF_REDIS_DATABASE) || undefined
    },
    migrate: process.env.RCTF_DATABASE_MIGRATE
  },
  instanceType: process.env.RCTF_INSTANCE_TYPE,
  tokenKey: process.env.RCTF_TOKEN_KEY,
  origin: process.env.RCTF_ORIGIN,
  ctftime: {
    clientId: process.env.RCTF_CTFTIME_CLIENT_ID,
    clientSecret: process.env.RCTF_CTFTIME_CLIENT_SECRET
  },
  userMembers: parseBoolEnv(process.env.RCTF_USER_MEMBERS),
  sponsors: undefined,
  homeContent: process.env.RCTF_HOME_CONTENT,
  ctfName: process.env.RCTF_NAME,
  meta: {
    description: process.env.RCTF_META_DESCRIPTION,
    imageUrl: process.env.RCTF_IMAGE_URL
  },
  logoUrl: process.env.RCTF_LOGO_URL,
  globalSiteTag: process.env.RCTF_GLOBAL_SITE_TAG,
  challengeProvider: {
    name: undefined,
    options: undefined
  },
  uploadProvider: {
    name: undefined,
    options: undefined
  },
  email: {
    provider: {
      name: undefined,
      options: undefined
    },
    from: process.env.RCTF_EMAIL_FROM
  },
  divisions: undefined,
  startTime: parseInt(process.env.RCTF_START_TIME) || undefined,
  endTime: parseInt(process.env.RCTF_END_TIME) || undefined,
  leaderboard: {
    maxLimit: parseInt(process.env.RCTF_LEADERBOARD_MAX_LIMIT) || undefined,
    maxOffset: parseInt(process.env.RCTF_LEADERBOARD_MAX_OFFSET) || undefined,
    updateInterval: parseInt(process.env.RCTF_LEADERBOARD_UPDATE_INTERVAL) || undefined,
    graphMaxTeams: parseInt(process.env.RCTF_LEADERBOARD_GRAPH_MAX_TEAMS) || undefined,
    graphSampleTime: parseInt(process.env.RCTF_LEADERBOARD_GRAPH_SAMPLE_TIME) || undefined
  },
  loginTimeout: parseInt(process.env.RCTF_LOGIN_TIMEOUT) || undefined
}

const defaultConfig: PartialDeep<ServerConfig> = {
  database: {
    migrate: 'never'
  },
  instanceType: 'all',
  userMembers: true,
  sponsors: [],
  homeContent: '',
  challengeProvider: {
    name: 'challenges/database'
  },
  uploadProvider: {
    name: 'uploads/local'
  },
  meta: {
    description: '',
    imageUrl: ''
  },
  leaderboard: {
    maxLimit: 100,
    maxOffset: 4294967296,
    updateInterval: 10000,
    graphMaxTeams: 10,
    graphSampleTime: 1800000
  },
  loginTimeout: 3600000
}

const _removeUndefined = <T>(o: Record<string, T>): Record<string, T> | undefined => {
  let hasKeys = false
  for (const key of Object.keys(o)) {
    let v = o[key]
    if (typeof v === 'object' && v != null) {
      o[key] = v = _removeUndefined(v as Record<string, unknown>) as T
    }
    if (v === undefined || v === null) {
      delete o[key]
    } else {
      hasKeys = true
    }
  }
  return hasKeys ? o : undefined
}

const removeUndefined = <T extends Record<string, unknown>>(o: T & Record<string, unknown>): T => {
  const cleaned = _removeUndefined(o) as T | undefined
  return cleaned ?? ({} as T)
}

const config = deepMerge.all([defaultConfig, ...fileConfigs, removeUndefined(envConfig)]) as ServerConfig

export default config
