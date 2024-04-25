import { DirectoryEntriesDataTable } from "@/components/tables/directory-entries-data-table";

export default async function DirectoryEntries({ params }: { params: { id: number } }) {
  return (
    <main className="container mx-auto py-10">

      <h1 className="scroll-m-20 text-4xl font-extrabold tracking-tight lg:text-5xl mb-4">
        Directory {params.id}
      </h1>
      <DirectoryEntriesDataTable id={params.id} />
    </main>
  );
}
