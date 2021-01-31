<script lang="ts">
	import { stats } from './store'
	import { createWebSocket } from './utils/ws'
	import BasicInfo from './components/BasicInfo.svelte'
	import Usage from './components/Usage.svelte'
	import Disk from './components/Disk.svelte'
	import Process from './components/Process.svelte'
	
</script>

<main class="container">
	<div class="app">
		<h1 class="title">Stop</h1>
		{#each $stats as s}
			{#if s}
				<div>
					<h2>{s.name}</h2>	
					<Usage bind:data={s} />
					<BasicInfo bind:data={s} />
					<Disk bind:data={s} />
					<Process bind:data={s} />
				</div>
			{/if}
		{/each}
	</div>
</main>

<style>
	.title {
		margin: 0;
	}

	.app {
		width: 680px;
		margin-left: auto;
		margin-right: auto;
	}
</style>