<script lang="ts">
	import { type Metadata, type MetadataElements, type MetadataType } from "$lib/api";

	let {
		elements,
		identifier = $bindable(),
		metadata = $bindable(),
	}: {
		elements: MetadataElements;
		identifier: string;
		metadata: Metadata[];
	} = $props();

	type FlatMeta = {
		type?: MetadataType;
		choices?: string[];
	};

	let flatMetadata = new Map<string, FlatMeta>();
	for (let { name, type, choices: c } of Object.values(elements).flat()) {
		let choices = c?.map(({ label }) => label);
		flatMetadata.set(name, { type, choices });
	}

	let meta = $state(
		elements.mandatory.map((kind) => ({
			mandatory: true,
			kind: kind.name,
			value: "",
		})),
	);

	$effect(() => {
		metadata = meta.map(({ kind, value }) => ({ kind, value }));
	});
</script>

<h2 class="text-lg">Metadata</h2>

<div class="table border-collapse p-2">
	<div class="mb-2 table-row">
		<span class="table-cell p-2 pl-0 font-bold">kind</span>
		<span class="table-cell p-2 font-bold">value</span>
	</div>

	<div class="table-row">
		<span class="table-cell pr-2">identifier</span>
		<span class="table-cell pl-2">
			<input class="h-6 w-full outline" placeholder="doi" bind:value={identifier} />
		</span>
	</div>

	{#each meta as m, i (m)}
		{@const { type, choices } = flatMetadata.get(m.kind)!}

		<div class="table-row">
			{#if m.mandatory}
				<span class="table-cell pr-2">{m.kind}</span>
			{:else}
				<select bind:value={m.kind} class="table-cell field-sizing-content pr-2">
					{#each Object.entries(elements) as [type, values]}
						<optgroup label={type}>
							{#each values as { name: kind }}
								<option value={kind}>{kind}</option>
							{/each}
						</optgroup>
					{/each}
				</select>
			{/if}

			<span class="table-cell pl-2">
				{#if type == "choice"}
					<select required class="h-6 w-full outline" bind:value={m.value}>
						{#each choices! as choice}
							<option value={choice}>{choice}</option>
						{/each}
					</select>
				{:else if type === "year"}
					<input
						placeholder={type}
						class="h-6 w-full outline"
						bind:value={
							() => m.value,
							(v) => {
								if (v === undefined) {
									m.value = "";
								} else if (!isNaN(+v)) {
									m.value = v;
								}
							}
						}
						type="text"
					/>
				{:else if type == "url"}
					<input placeholder={type} class="h-6 w-full outline" bind:value={m.value} type="url" />
				{:else}
					<input placeholder={type || m.kind} class="h-6 w-full outline" bind:value={m.value} />
				{/if}
			</span>
			{#if !m.mandatory}
				<button
					onclick={() => meta.splice(i, 1)}
					class="cursor-pointer pl-2 font-bold text-red-500"
				>
					x
				</button>
			{:else}
				<button class="invisible pl-2 font-bold">x</button>
			{/if}
		</div>
	{/each}

	<div class="table-caption caption-bottom">
		<div class="flex w-full justify-center">
			<button
				onclick={() => meta.push({ mandatory: false, kind: "title", value: "" })}
				class="cursor-pointer px-4 py-2 text-center font-bold">+ add metadata</button
			>
		</div>
	</div>
</div>
