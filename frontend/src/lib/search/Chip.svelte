<script lang="ts" module>
	export type ChipElement = {
		kind: string;
		value?: string;
	};
</script>

<script lang="ts">
	let { chip } = $props();
</script>

<span class="flex items-center gap-0.5">
	<span class="flex h-8 items-center rounded-l bg-neutral-200 px-2 text-nowrap dark:bg-neutral-600">
		{chip.kind}
	</span>
	<span class="flex h-8 items-center bg-neutral-200 px-2 dark:bg-neutral-600">≈</span>

	{#if chip.value}
		<span
			class="flex h-8 items-center rounded-r bg-neutral-200 px-2 text-nowrap dark:bg-neutral-600"
		>
			{chip.value}
		</span>
	{/if}
</span>
