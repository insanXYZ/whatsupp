import { useMutation, useQuery } from "@tanstack/react-query";
import { API } from "./axios";

function Query(queryKey: any[], url: string) {
  useQuery({
    queryKey,
    queryFn: async () => {
      const { data } = await API.get(url);
      return data;
    },
  });
}

export enum HttpMethod {
  POST = "POST",
  PUT = "PUT",
  DELETE = "DELETE"
}

interface Mutate {
  url: string
  body: any
  method: HttpMethod
}

function Mutation(mutationKey: any[]) {
  return useMutation({
    mutationFn: async ({ url, body, method }: Mutate) => {
      const res = await API({
        url: url,
        data: body,
        method: method
      })

      return res.data
    },
    mutationKey,
  });
}

export { Query, Mutation };
