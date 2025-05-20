<script lang="ts" module>
	export type Node = FlowNode<Omit<ApiNode, 'position'>, 'custom'>;
</script>

<script lang="ts">
	import { type NodeProps, type Node as FlowNode, Position } from '@xyflow/svelte';
	import type { Node as ApiNode } from '$lib/api';
	import { Handle } from '@xyflow/svelte';

	let { data }: NodeProps<Node> = $props();

	function isLight(s: string) {
		const rgb = parseInt(s.slice(1), 16);
		const r = (rgb >> 16) & 0xff;
		const g = (rgb >> 8) & 0xff;
		const b = rgb & 0xff;

		const luma = r * 0.2126 + g * 0.7152 + b * 0.0722;
		return luma > 127;
	}
</script>

<div
	class="rounded-lg px-4 py-2"
	class:text-white={!isLight(data.color)}
	style:background-color={data.color}
>
	{data.body}
</div>

<Handle type="target" position={Position.Left} />
<Handle type="source" position={Position.Right} />
