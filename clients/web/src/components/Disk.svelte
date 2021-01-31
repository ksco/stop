<script lang="ts">
  import type { Stats } from '../store'

  interface Props extends Pick<Stats, 'disks' | 'diskTotal' | 'diskUsage'> {}
  export const data:Props | null = null
</script>
<style></style>