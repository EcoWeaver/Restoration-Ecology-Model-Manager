import type { Node as FlowNode, Edge as FlowEdge } from "@xyflow/svelte";
import { putModel, type Node as ApiNode, type Relation as ApiRelation, type Metadata, type ModelRequest } from "./api"

export type Node = FlowNode<Omit<ApiNode, "position"> & { editable: boolean }, "custom">;
export type Edge = FlowEdge<Omit<ApiRelation, "from" | "to"> & { editable: boolean }, "custom">;

export async function saveModel(identifier: string, metadata: Metadata[], n: Node[], e: Edge[], version?: string): Promise<ModelRequest> {
	const nodes: ApiNode[] = n.map((n) => {
		let { body, color } = n.data;
		return {
			body,
			color,
			position: {
				x: Math.round(n.position.x),
				y: Math.round(n.position.y),
			},
		};
	});

	const relations: ApiRelation[] = e
		.filter((e) => e.type == "custom")
		.map((e) => {
			const from = n.findIndex((n) => n.id === e.source);
			const to = n.findIndex((n) => n.id === e.target)!;

			let { kind } = e.data!;
			return { kind, from, to };
		});


	const graph = { nodes, relations };
	const request: ModelRequest = {
		identifier,
		version: version ?? "",
		metadata,
		graph
	};

	return await putModel(request);
}
