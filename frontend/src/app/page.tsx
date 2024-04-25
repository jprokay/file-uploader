import GetStartedAlert from "@/components/alerts/get-started-alert";
import { DirectoriesDataTable } from "@/components/tables/directories-data-table";
import { DirectoriesWithTotal, client } from "@/lib/contacts-api";
import Head from "next/head";
import { cookies } from "next/headers";

type GetDataProps = {
  limit: number
  offset: number
  sort: "asc" | "desc" | undefined
}
async function getData(props: GetDataProps): Promise<DirectoriesWithTotal> {
  const res = await client.GET("/directories", {
    params: {
      query: { sort: props.sort, limit: props.limit, offset: props.offset },
      cookie: { userId: cookies().get("userId")?.value || "" }
    },
    headers: { Cookie: cookies().toString() },
  })


  if (res.error) {
    return { directories: [], total: 0 }
  }

  return res.data
}


export default async function Home() {
  const data = await getData({ limit: 10, offset: 1, sort: "asc" })
  return (
    <div>
      <main className="container mx-auto py-10">
        <h1 className="scroll-m-20 text-4xl font-extrabold tracking-tight lg:text-5xl mb-4">
          Directories
        </h1>
        {data.total < 1 && <div className="my-4">
          <GetStartedAlert />
        </div>}
        <DirectoriesDataTable rowCount={data.total} data={data.directories} />
      </main>
    </div>
  );
}
