package resolvers

import (
	awsContext "context"
	"encoding/json"
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Repository", func() {
	type arguments struct {
		Bar string `json:"bar"`
	}
	type response struct {
		Foo string
	}

	r := New()
	resolver := func(ctx awsContext.Context, arg arguments) (response, error) { return response{"bar"}, nil }
	resolverWithError := func(ctx awsContext.Context, arg arguments) (response, error) {
		return response{"bar"}, errors.New("Has Error")
	}

	_ = r.Add("example.resolver", resolver)
	_ = r.Add("example.resolver.with.error", resolverWithError)

	Context("Matching invocation", func() {
		res, err := r.Handle(awsContext.TODO(), invocation{
			Resolve: "example.resolver",
			Context: context{
				Arguments: json.RawMessage(`{"bar":"foo"}`),
			},
		})

		It("Should not error", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		It("Should have data", func() {
			Expect(res.(response).Foo).To(Equal("bar"))
		})
	})

	Context("Matching invocation with error", func() {
		_, err := r.Handle(awsContext.TODO(), invocation{
			Resolve: "example.resolver.with.error",
			Context: context{
				Arguments: json.RawMessage(`{"bar":"foo"}`),
			},
		})

		It("Should error", func() {
			Expect(err).To(HaveOccurred())
		})
	})

	Context("Matching invocation with invalid payload", func() {
		_, err := r.Handle(awsContext.TODO(), invocation{
			Resolve: "example.resolver.with.error",
			Context: context{
				Arguments: json.RawMessage(`{"bar:foo"}`),
			},
		})

		It("Should error", func() {
			Expect(err).To(HaveOccurred())
		})
	})

	Context("Not matching invocation", func() {
		res, err := r.Handle(awsContext.TODO(), invocation{
			Resolve: "example.resolver.not.found",
			Context: context{
				Arguments: json.RawMessage(`{"bar":"foo"}`),
			},
		})

		It("Should error", func() {
			Expect(err).To(HaveOccurred())
		})

		It("Should have no data", func() {
			Expect(res).To(BeNil())
		})
	})
})
