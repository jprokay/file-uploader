import { paths } from "./contacts";
import createClient from "openapi-fetch";

export const client = createClient<paths>({
	baseUrl: process.env.NEXT_PUBLIC_BACKEND_URL || process.env.BACKEND_URL,
	credentials: "include",
});

export type ContactsWithTotal =
	paths["/contacts"]["get"]["responses"]["200"]["content"]["application/json"];
export type Contacts = ContactsWithTotal["items"];
export type Contact = Contacts[0];

export type DirectoriesWithTotal =
	paths["/directories"]["get"]["responses"]["200"]["content"]["application/json"];

export type Directories = DirectoriesWithTotal["items"];
export type Directory = Directories[0];

export type DirectoryEntriesWithTotal =
	paths["/directories/{id}/entries"]["get"]["responses"]["200"]["content"]["application/json"];

export type DirectoryEntries = DirectoryEntriesWithTotal["items"];
export type DirectoryEntry = DirectoryEntries[0];
