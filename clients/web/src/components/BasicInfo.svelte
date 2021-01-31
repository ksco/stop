<script lang="ts">
  import type { Stats } from '../store'

  interface Props extends Pick<Stats, 'osName' | 'addr' | 'uptime' | 'cpuCores'> {}

  export let data: Props | null = null
</script>
<div>
  <header class="header">
    <h3 class="title">
      {data.osName}
    </h3>
    <div class="ip">{data.addr}</div>
  </header>
  <div class="body">
    <div class="item">
      <div class="label">uptime</div>
      <div class="value">{data.uptime}</div>
    </div>
    <div class="item">
      <div class="label">cores</div>
      <div class="value">{data.cpuCores}</div>
    </div>
  </div>
</div>
<style>
  .header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .title {
    font-size: 16px;
    margin: 0;
  }

  .body {
    display: flex;
    flex-wrap: wrap;
  }

  .item {
    flex: 0 0 25%;
    max-width: 25%;
    padding: 10px;
  }
</style>