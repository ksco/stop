import { writable } from 'svelte/store'
import data from '../mock/data.json'

interface Process {
 user: string;
 cpuPercent: number;
 rss: number; 
 cmd: string; 
}

interface Disk {
  fileSystem: string; 
  total: number;
  usage: number;
}

interface Load {
  one: number;
  five: number;
  fifteen: number;
}

export interface Stats {
  name: string; 
  addr: string;
  uptime: number;
  sessionsCount: number;
  processesCount: number;
  processes: Process[];
  fileHandlesCount: number;
  fileHandlesLimit: number;
  osKernel: string;
  osName: string;
  osArch: string;
  cpuName: string;
  cpuCores: number;
  cpuFreq: number;
  cpuUsagePercent: number;
  ramTotal: number;
  ramUsage: number;
  swapTotal: number;
  swapUsage: number;
  disks: Disk[];
  diskTotal: number;
  diskUsage: number;
  connectionsCount: number;
  load: Load
}

function createStats () {
  const { subscribe, set, update } = writable<Stats[]>([data])

  return { 
    subscribe,
    set,
    update,
    reset: () => set(null)
  }
}

export const stats = createStats()