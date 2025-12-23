<script lang="ts">
	import { apiToFlow, flowToApi } from "$lib/model.js";
	import GraphEdit from "$lib/GraphEdit.svelte";
	import { SvelteFlowProvider } from "@xyflow/svelte";
	import { page } from "$app/state";
	import { deleteLocalModel, setLocalModel } from "$lib/local.js";
	import { CheckLine, Upload } from "@lucide/svelte";
	import { clearNavigate, startNavigate } from "$lib/LoadingSpinner.svelte";
	import { req, type BackendError, type PutModelResult } from "$lib/api";
	import { resolve } from "$app/paths";
	import { goto } from "$app/navigation";
	import Alert from "$lib/Alert.svelte";

	let { data } = $props();
	const { metadataElements, model } = data;

	const graph = apiToFlow(model.graph, true);
	const parent = model.parent;

	let nodes = $state.raw(graph.nodes);
	let edges = $state.raw(graph.edges);
	let metadata = $state(model.metadata);

	$effect(() => {
		const nodes_ = nodes;
		const edges_ = edges;
		const metadata_ = metadata;
		setLocalModel(page.params.model!, {
			graph: flowToApi(nodes_, edges_),
			metadata: metadata_,
			parent,
		});
	});

	const title = $derived(metadata.find(({ kind }) => kind === "title")?.value ?? "");

	let alert: ReturnType<typeof Alert>;

	async function check() {
		try {
			startNavigate();
			await req("POST", "model/verify", {
				graph: flowToApi(nodes, edges),
				metadata,
				parent: parent?.uuid,
			});
			alert.info("model OK");
			clearNavigate();
		} catch (err: unknown) {
			const error = err as BackendError;
			alert.error(error.message);
			clearNavigate();
		}
	}

	async function upload() {
		try {
			startNavigate();
			const res = await req<PutModelResult>("PUT", "model", {
				graph: flowToApi(nodes, edges),
				metadata,
				parent: parent?.uuid,
			});
			deleteLocalModel(page.params.model!);
			await goto(resolve(`/view/${res.uuid}`));
		} catch (err: unknown) {
			const error = err as BackendError;
			alert.error(error.message);
			clearNavigate();
		}
	}
</script>

<svelte:head>
	<title>edit '{title}'</title>
</svelte:head>

<Alert bind:this={alert} />

<SvelteFlowProvider>
	<GraphEdit {parent} {metadataElements} bind:metadata bind:nodes bind:edges>
		<button
			class="flex h-10 cursor-pointer items-center rounded-md bg-indigo-300 px-3 outline-offset-2 outline-indigo-400 hover:bg-indigo-400 focus-visible:outline-2 dark:bg-indigo-500"
			title="check"
			onclick={check}
		>
			<CheckLine />
			<span class="pl-2">check</span>
		</button>
		<button
			class="flex h-10 cursor-pointer items-center rounded-md bg-indigo-300 px-3 outline-offset-2 outline-indigo-400 hover:bg-indigo-400 focus-visible:outline-2 dark:bg-indigo-500"
			title="upload"
			onclick={upload}
		>
			<Upload />
			<span class="pl-2">upload</span>
		</button>
	</GraphEdit>
</SvelteFlowProvider>
