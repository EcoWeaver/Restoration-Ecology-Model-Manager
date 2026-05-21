<script lang="ts">
	import { Link, PencilLine } from "@lucide/svelte";
	import Graph from "$lib/Graph.svelte";
	import PanelComponent from "$lib/graph/PanelComponent.svelte";
	import Header from "$lib/Header.svelte";
	import Alert from "$lib/Alert.svelte";
	import { apiToFlow } from "$lib/model";
	import { SvelteFlowProvider } from "@xyflow/svelte";
	import { resolve } from "$app/paths";
	import { newLocalModel } from "$lib/local.js";
	import { goto } from "$app/navigation";
	import { useElements, type MetadataType } from "$lib/api.js";

	let { data } = $props();

	const { getMetadataElements } = useElements();

	let { model } = $derived(data);
	const { uuid } = $derived(model);
	let { graph, metadata } = $derived(model.model);
	let { nodes, edges } = $derived(apiToFlow(graph, false));
	const title = $derived(metadata.find(({ kind }) => kind === "title")!.value);

	const flatMetadata = new Map<string, MetadataType | undefined>();
	for (const { name, type } of Object.values(getMetadataElements()).flat()) {
		flatMetadata.set(name, type);
	}

	let alert: ReturnType<typeof Alert>;

	async function edit() {
		const localUUID = newLocalModel({
			graph,
			metadata: [
				{ kind: "derived from", value: uuid },
				...metadata.filter((m) => m.kind !== "derived from"),
			],
		});

		await goto(resolve(`/edit/${localUUID}`));
	}

	async function share() {
		await navigator.clipboard.writeText(window.location.href);
		alert.info("link copied to clipboard");
	}
</script>

<svelte:head>
	<title>{title}</title>
</svelte:head>

<div class="absolute z-10 w-full bg-white dark:bg-neutral-800">
	<Header>
		<button
			class="flex h-10 cursor-pointer items-center rounded-md bg-indigo-300 px-3 outline-offset-2 outline-indigo-400 hover:bg-indigo-400 focus-visible:outline-2 dark:bg-indigo-500"
			title="edit"
			onclick={edit}
		>
			<PencilLine strokeWidth={1.75} />
			<div class="pl-2">edit</div>
		</button>

		<button
			class="flex h-10 cursor-pointer items-center rounded-md bg-indigo-300 px-3 outline-offset-2 outline-indigo-400 hover:bg-indigo-400 focus-visible:outline-2 dark:bg-indigo-500"
			title="copy model link to clipboard"
			onclick={share}
		>
			<Link />
			<div class="pl-2">share</div>
		</button>
	</Header>
</div>

<Alert bind:this={alert} />

<div class="h-dvh w-dvw">
	<SvelteFlowProvider>
		<Graph {nodes} {edges}>
			<PanelComponent>
				<div class="table">
					<div class="table-row">
						<h2 class="table-cell pr-2 pb-2 text-xl font-bold select-text">Metadata</h2>
					</div>
					<div class="table-row">
						<div class="table-cell py-1 pr-2 font-bold select-text">kind</div>
						<div class="table-cell py-1 pl-2 font-bold select-text">value</div>
					</div>
					{#each metadata as { kind, value }, idx (idx)}
						{@const type = flatMetadata.get(kind)}

						<div class="table-row">
							<div class="table-cell py-1 pr-2 select-text">{kind}</div>

							{#if type === "doi"}
								<a
									class="table-cell max-w-lg py-1 pl-2 text-blue-700 outline-indigo-400 select-text hover:underline focus-visible:outline-2 dark:text-blue-400"
									href="https://doi.org/{value}"
									target="_blank"
									rel="external noreferrer"
								>
									{value}
								</a>
							{:else if type === "orcid"}
								<a
									class="table-cell max-w-lg py-1 pl-2 text-blue-700 outline-indigo-400 select-text hover:underline focus-visible:outline-2 dark:text-blue-400"
									href="https://orcid.org/{value}"
									target="_blank"
									rel="external noreferrer"
								>
									{value}
								</a>
							{:else if type === "url"}
								<a
									class="table-cell max-w-lg py-1 pl-2 text-blue-700 outline-indigo-400 select-text hover:underline focus-visible:outline-2 dark:text-blue-400"
									href={value}
									target="_blank"
									rel="external noreferrer"
								>
									{value}
								</a>
							{:else if kind === "derived from"}
								<a
									class="table-cell max-w-lg py-1 pl-2 text-blue-700 outline-indigo-400 select-text hover:underline focus-visible:outline-2 dark:text-blue-400"
									href={resolve(`/view/${value}`)}
								>
									{value}
								</a>
							{:else}
								<div class="table-cell max-w-lg py-1 pl-2 select-text">{value}</div>
							{/if}
						</div>
					{/each}
				</div>
			</PanelComponent>
		</Graph>
	</SvelteFlowProvider>
</div>
