import type { Graph, Metadata, PostSearchResult } from "./api";

export type LocalModel = {
	graph: Graph;
	metadata: Metadata[];
};

export function getLocalModels(): Map<string, Omit<PostSearchResult, "uuid">> {
	return new Map(Object.entries(JSON.parse(localStorage.getItem("localModels") ?? "{}")));
}

export function newLocalModel(model?: LocalModel): string {
	const localUUID = self.crypto.randomUUID();
	setLocalModel(localUUID, model ?? { graph: {}, metadata: [] });
	return localUUID;
}

export function getLocalModel(localUUID: string): LocalModel | undefined {
	const localModel = localStorage.getItem(localUUID);
	return localModel ? JSON.parse(localModel) : undefined;
}

export function setLocalModel(localUUID: string, model: LocalModel) {
	withLocalModels((localModels) => localModels.set(localUUID, searchResult(model)));
	localStorage.setItem(localUUID, JSON.stringify(model));
}

export function deleteLocalModel(localUUID: string) {
	withLocalModels((localModels) => localModels.delete(localUUID));
	localStorage.removeItem(localUUID);
}

function withLocalModels(f: (localModels: Map<string, Omit<PostSearchResult, "uuid">>) => void) {
	const localModels = getLocalModels();
	f(localModels);
	localStorage.setItem("localModels", JSON.stringify(Object.fromEntries(localModels)));
}

function searchResult(model: LocalModel): Omit<PostSearchResult, "uuid"> {
	return {
		uploaded: new Date().toISOString(),
		identifier: model.metadata.find(({ kind }) => kind === "identifier")?.value ?? "",
		title: model.metadata.find(({ kind }) => kind === "title")?.value ?? "",
		modelType: model.metadata.find(({ kind }) => kind === "model type")?.value ?? "",
	};
}
