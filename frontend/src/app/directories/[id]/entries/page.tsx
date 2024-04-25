import { DirectoryEntriesDataTable } from "@/components/tables/directory-entries-data-table";
import { DirectoryEntriesWithTotal, client } from "@/lib/contacts-api";
import { cookies } from "next/headers";

type GetDataProps = {
  limit: number
  offset: number
}
async function getData(id: number, props: GetDataProps): Promise<DirectoryEntriesWithTotal> {
  const res = await client.GET("/directories/{id}/entries", {
    params: {
      query: props,
      cookie: { userId: cookies().get("userId")?.value || "" },
      path: {
        id
      }
    },
    headers: { Cookie: cookies().toString() },
  })


  if (res.error) {
    return { entries: [], total: 0 }
  }

  return res.data
}


export default async function DirectoryEntries({ params }: { params: { id: number } }) {
  const data = await getData(params.id, { limit: 100, offset: 0 })

  console.log("Entries Data: ", data)
  return (
    <main className="container mx-auto py-10">

      <h1 className="scroll-m-20 text-4xl font-extrabold tracking-tight lg:text-5xl mb-4">
        Directory {params.id}
      </h1>
      <DirectoryEntriesDataTable id={params.id} data={data.entries} rowCount={data.total} />
    </main>
  );
}
