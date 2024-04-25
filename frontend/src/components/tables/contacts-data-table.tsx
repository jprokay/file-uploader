"use client"
import { Contact, client } from "@/lib/contacts-api"
import { DataTable } from "./data-table"
import { ColumnDef, PaginationState } from "@tanstack/react-table"
import { cookies, getCookieValue } from "@/lib/cookies"

const columns: ColumnDef<Contact>[] = [
	{
		accessorKey: "contact_id",
		header: "ID",
	},
	{
		accessorKey: "contact_first_name",
		header: "First Name",
	},
	{
		accessorKey: "contact_last_name",
		header: "Last Name",
	},
	{
		accessorKey: "contact_email",
		header: "Email",
	},
];

async function getData(props: PaginationState) {

	const res = await client.GET("/contacts", {
		params: {
			query: { limit: props.pageSize, offset: props.pageIndex },
			cookie: { userId: getCookieValue("userId") || "" }
		},
		headers: { Cookie: cookies() },
	})


	if (res.error) {
		throw new Error('Oops something went wrong')
	}

	return res.data
}


export const ContactsDataTable = () => {
	return <DataTable
		queryKey={'contacts'} defaultPageSize={100} pageSizes={[25, 50, 100, 250, 500]}
		columns={columns}
		queryFn={(state) => () => getData(state)} />
}
