<script lang="ts">
  import type { Stats } from '../store'

  interface Props extends Pick<Stats, 'ramTotal' | 'ramUsage' | 'cpuUsagePercent'> {}
  
  export const data: Props | null = null
</script>

<div class="block">
  <div class="chart">{data?.ramTotal}</div>
</div>
<style></style>