<script lang="ts">
	import { apiToFlow, flowToApi } from "$lib/model.js";
	import GraphEdit from "$lib/GraphEdit.svelte";
	import { SvelteFlowProvider } from "@xyflow/svelte";
	import { page } from "$app/state";
	import { deleteLocalModel, setLocalModel } from "$lib/local.js";
	import { Check, CheckLine, Merge, Plus, Upload, X } from "@lucide/svelte";
	import { clearNavigate, startNavigate } from "$lib/LoadingSpinner.svelte";
	import { req, type BackendError, type GetModelResult, type PutModelResult } from "$lib/api";
	import { resolve } from "$app/paths";
	import { beforeNavigate, goto } from "$app/navigation";
	import Alert from "$lib/Alert.svelte";
	import { untrack } from "svelte";
	import UuidInput from "$lib/metadata-panel/UuidInput.svelte";

	let { data } = $props();

	const { model } = $derived(data);
	let { nodes, edges } = $derived(apiToFlow(model.graph, true));
	let metadata = $derived(model.metadata);
	const title = $derived(metadata.find(({ kind }) => kind === "title")?.value ?? "");

	let timeoutId = $state<NodeJS.Timeout>();
	$effect(() => {
		const nodes_ = nodes;
		const edges_ = edges;
		const metadata_ = metadata;

		clearTimeout(untrack(() => timeoutId));
		timeoutId = setTimeout(() => {
			setLocalModel(page.params.model!, {
				graph: flowToApi(nodes_, edges_),
				metadata: metadata_,
			});
		}, 50);
	});

	beforeNavigate((nav) => {
		if (nav.type !== "goto" || nav.to?.route.id !== "/view/[...model]") {
			if (localStorage.getItem("editModelConfirmNavigation")) {
				return;
			}

			const doNavigate = confirm(
				"Leave site?\n\nYour local models are saved in your browser's local storage until they're uploaded. They'll be deleted if this page's data is cleared.",
			);

			if (doNavigate) {
				localStorage.setItem("editModelConfirmNavigation", "true");
			} else {
				nav.cancel();
			}
		}
	});

	let alert: ReturnType<typeof Alert>;

	async function check() {
		try {
			startNavigate();
			await req("POST", "model/verify", {
				model: {
					graph: flowToApi(nodes, edges),
					metadata,
				},
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
				model: { graph: flowToApi(nodes, edges), metadata },
			});
			deleteLocalModel(page.params.model!);
			await goto(resolve(`/view/${res.uuid}`));
		} catch (err: unknown) {
			const error = err as BackendError;
			alert.error(error.message);
			clearNavigate();
		}
	}

	let mergeUuid = $state("");
	let mergeModels = $state<GetModelResult[]>([]);

	async function pushMerge() {
		try {
			if (mergeUuid.length !== 36) {
				return;
			}
			startNavigate();
			mergeModels.push(await req<GetModelResult>("GET", `model?uuid=${mergeUuid}`));
			mergeUuid = "";
		} catch (err: unknown) {
			const error = err as BackendError;
			alert.error(error.message);
			clearNavigate();
		}
	}

	let graphEdit: ReturnType<typeof GraphEdit>;
	let dialogElement = $state<HTMLDialogElement>();

	let mergedModels = 0;
	async function merge() {
		await pushMerge();
		for (const model of mergeModels) {
			graphEdit.addMetadata({ kind: "derived from", value: model.uuid });
			const graph = apiToFlow(model.model.graph, true, `${mergedModels++}_`);
			nodes = [...nodes, ...graph.nodes];
			edges = [...edges, ...graph.edges];
		}
		mergeModels = [];
		dialogElement?.close();
	}
</script>

<svelte:head>
	<title>edit '{title}'</title>
</svelte:head>

<Alert bind:this={alert} />

<SvelteFlowProvider>
	<GraphEdit bind:this={graphEdit} bind:metadata bind:nodes bind:edges>
		<button
			class="flex h-10 cursor-pointer items-center rounded-md bg-indigo-300 px-3 outline-offset-2 outline-indigo-400 hover:bg-indigo-400 focus-visible:outline-2 dark:bg-indigo-500"
			title="check"
			onclick={check}
		>
			<CheckLine />
			<div class="pl-2">check</div>
		</button>
		<button
			class="flex h-10 cursor-pointer items-center rounded-md bg-indigo-300 px-3 outline-offset-2 outline-indigo-400 hover:bg-indigo-400 focus-visible:outline-2 dark:bg-indigo-500"
			title="upload"
			onclick={upload}
		>
			<Upload />
			<div class="pl-2">upload</div>
		</button>
		<button
			class="cursor-pointer rounded-full outline-offset-4 outline-indigo-400 focus-visible:outline-2"
			title="merge models"
			onclick={() => dialogElement?.showModal()}
		>
			<Merge />
		</button>
	</GraphEdit>
</SvelteFlowProvider>

<dialog
	bind:this={dialogElement}
	onclose={() => (mergeUuid = "")}
	closedby="any"
	class="m-auto rounded-md bg-neutral-200 p-2 outline-2 outline-indigo-400 backdrop:bg-white/50 dark:bg-neutral-700 dark:backdrop:bg-neutral-800/50"
>
	<div class="mb-2 flex items-center">
		Enter the IDs of the to-be-merged models.
		<button
			class="ml-1 flex cursor-pointer items-center rounded-full p-1 outline-offset-2 outline-neutral-500 hover:bg-neutral-500/25 focus-visible:outline-2"
			onclick={merge}
		>
			<Check />
		</button>
	</div>
	<div class="flex items-center">
		<UuidInput bind:value={mergeUuid} />
		<button
			class="ml-1 flex cursor-pointer items-center rounded-full p-1 outline-offset-2 outline-neutral-500 hover:bg-neutral-500/25 focus-visible:outline-2"
			onclick={pushMerge}
		>
			<Plus />
		</button>
	</div>
	{#each mergeModels as mergeModel, idx (idx)}
		<div class="mt-1 flex items-center pr-1 pl-3.25">
			<a
				class="w-full text-blue-700 outline-indigo-400 select-text hover:underline focus-visible:outline-2 dark:text-blue-400"
				target="_blank"
				href={resolve(`/view/${mergeModel.uuid}`)}
			>
				{mergeModel.uuid}
			</a>
			<button
				class="float-right h-6 cursor-pointer rounded-full outline-red-500 hover:bg-red-500/25 focus-visible:outline-2"
				title="delete this metadata entry"
				onclick={() => mergeModels.splice(idx, 1)}
			>
				<X class="stroke-red-500" />
			</button>
		</div>
	{/each}
</dialog>
