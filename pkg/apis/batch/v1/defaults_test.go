/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1_test

import (
	"reflect"
	"testing"

	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubernetes/pkg/api/legacyscheme"
	_ "k8s.io/kubernetes/pkg/apis/batch/install"
	_ "k8s.io/kubernetes/pkg/apis/core/install"
	utilpointer "k8s.io/utils/pointer"

	. "k8s.io/kubernetes/pkg/apis/batch/v1"
)

func TestSetDefaultJob(t *testing.T) {
	defaultLabels := map[string]string{"default": "default"}
	tests := map[string]struct {
		original     *batchv1.Job
		expected     *batchv1.Job
		expectLabels bool
	}{
		"All unspecified -> sets all to default values": {
			original: &batchv1.Job{
				Spec: batchv1.JobSpec{
					Template: v1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{Labels: defaultLabels},
					},
				},
			},
			expected: &batchv1.Job{
				Spec: batchv1.JobSpec{
					Completions:    utilpointer.Int32Ptr(1),
					Parallelism:    utilpointer.Int32Ptr(1),
					BackoffLimit:   utilpointer.Int32Ptr(6),
					CompletionMode: batchv1.NonIndexedCompletion,
				},
			},
			expectLabels: true,
		},
		"All unspecified -> all integers are defaulted and no default labels": {
			original: &batchv1.Job{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"mylabel": "myvalue"},
				},
				Spec: batchv1.JobSpec{
					Template: v1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{Labels: defaultLabels},
					},
				},
			},
			expected: &batchv1.Job{
				Spec: batchv1.JobSpec{
					Completions:    utilpointer.Int32Ptr(1),
					Parallelism:    utilpointer.Int32Ptr(1),
					BackoffLimit:   utilpointer.Int32Ptr(6),
					CompletionMode: batchv1.NonIndexedCompletion,
				},
			},
		},
		"WQ: Parallelism explicitly 0 and completions unset -> BackoffLimit is defaulted": {
			original: &batchv1.Job{
				Spec: batchv1.JobSpec{
					Parallelism: utilpointer.Int32Ptr(0),
					Template: v1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{Labels: defaultLabels},
					},
				},
			},
			expected: &batchv1.Job{
				Spec: batchv1.JobSpec{
					Parallelism:    utilpointer.Int32Ptr(0),
					BackoffLimit:   utilpointer.Int32Ptr(6),
					CompletionMode: batchv1.NonIndexedCompletion,
				},
			},
			expectLabels: true,
		},
		"WQ: Parallelism explicitly 2 and completions unset -> BackoffLimit is defaulted": {
			original: &batchv1.Job{
				Spec: batchv1.JobSpec{
					Parallelism: utilpointer.Int32Ptr(2),
					Template: v1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{Labels: defaultLabels},
					},
				},
			},
			expected: &batchv1.Job{
				Spec: batchv1.JobSpec{
					Parallelism:    utilpointer.Int32Ptr(2),
					BackoffLimit:   utilpointer.Int32Ptr(6),
					CompletionMode: batchv1.NonIndexedCompletion,
				},
			},
			expectLabels: true,
		},
		"Completions explicitly 2 and others unset -> parallelism and BackoffLimit are defaulted": {
			original: &batchv1.Job{
				Spec: batchv1.JobSpec{
					Completions: utilpointer.Int32Ptr(2),
					Template: v1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{Labels: defaultLabels},
					},
				},
			},
			expected: &batchv1.Job{
				Spec: batchv1.JobSpec{
					Completions:    utilpointer.Int32Ptr(2),
					Parallelism:    utilpointer.Int32Ptr(1),
					BackoffLimit:   utilpointer.Int32Ptr(6),
					CompletionMode: batchv1.NonIndexedCompletion,
				},
			},
			expectLabels: true,
		},
		"BackoffLimit explicitly 5 and others unset -> parallelism and completions are defaulted": {
			original: &batchv1.Job{
				Spec: batchv1.JobSpec{
					BackoffLimit: utilpointer.Int32Ptr(5),
					Template: v1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{Labels: defaultLabels},
					},
				},
			},
			expected: &batchv1.Job{
				Spec: batchv1.JobSpec{
					Completions:    utilpointer.Int32Ptr(1),
					Parallelism:    utilpointer.Int32Ptr(1),
					BackoffLimit:   utilpointer.Int32Ptr(5),
					CompletionMode: batchv1.NonIndexedCompletion,
				},
			},
			expectLabels: true,
		},
		"All set -> no change": {
			original: &batchv1.Job{
				Spec: batchv1.JobSpec{
					Completions:    utilpointer.Int32Ptr(8),
					Parallelism:    utilpointer.Int32Ptr(9),
					BackoffLimit:   utilpointer.Int32Ptr(10),
					CompletionMode: batchv1.NonIndexedCompletion,
					Template: v1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{Labels: defaultLabels},
					},
				},
			},
			expected: &batchv1.Job{
				Spec: batchv1.JobSpec{
					Completions:    utilpointer.Int32Ptr(8),
					Parallelism:    utilpointer.Int32Ptr(9),
					BackoffLimit:   utilpointer.Int32Ptr(10),
					CompletionMode: batchv1.NonIndexedCompletion,
					Template: v1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{Labels: defaultLabels},
					},
				},
			},
			expectLabels: true,
		},
		"All set, flipped -> no change": {
			original: &batchv1.Job{
				Spec: batchv1.JobSpec{
					Completions:    utilpointer.Int32Ptr(11),
					Parallelism:    utilpointer.Int32Ptr(10),
					BackoffLimit:   utilpointer.Int32Ptr(9),
					CompletionMode: batchv1.IndexedCompletion,
					Template: v1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{Labels: defaultLabels},
					},
				},
			},
			expected: &batchv1.Job{
				Spec: batchv1.JobSpec{
					Completions:    utilpointer.Int32Ptr(11),
					Parallelism:    utilpointer.Int32Ptr(10),
					BackoffLimit:   utilpointer.Int32Ptr(9),
					CompletionMode: batchv1.IndexedCompletion,
				},
			},
			expectLabels: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {

			original := test.original
			expected := test.expected
			obj2 := roundTrip(t, runtime.Object(original))
			actual, ok := obj2.(*batchv1.Job)
			if !ok {
				t.Fatalf("Unexpected object: %v", actual)
			}

			validateDefaultInt32(t, "Completions", actual.Spec.Completions, expected.Spec.Completions)
			validateDefaultInt32(t, "Parallelism", actual.Spec.Parallelism, expected.Spec.Parallelism)
			validateDefaultInt32(t, "BackoffLimit", actual.Spec.BackoffLimit, expected.Spec.BackoffLimit)

			if test.expectLabels != reflect.DeepEqual(actual.Labels, actual.Spec.Template.Labels) {
				if test.expectLabels {
					t.Errorf("Expected labels: %v, got: %v", actual.Spec.Template.Labels, actual.Labels)
				} else {
					t.Errorf("Unexpected equality: %v", actual.Labels)
				}
			}
			if actual.Spec.CompletionMode != expected.Spec.CompletionMode {
				t.Errorf("Got CompletionMode: %v, want: %v", actual.Spec.CompletionMode, expected.Spec.CompletionMode)
			}
		})
	}
}

func validateDefaultInt32(t *testing.T, field string, actual *int32, expected *int32) {
	if (actual == nil) != (expected == nil) {
		t.Errorf("Got different *%s than expected: %v %v", field, actual, expected)
	}
	if actual != nil && expected != nil {
		if *actual != *expected {
			t.Errorf("Got different %s than expected: %d %d", field, *actual, *expected)
		}
	}
}

func roundTrip(t *testing.T, obj runtime.Object) runtime.Object {
	data, err := runtime.Encode(legacyscheme.Codecs.LegacyCodec(SchemeGroupVersion), obj)
	if err != nil {
		t.Errorf("%v\n %#v", err, obj)
		return nil
	}
	obj2, err := runtime.Decode(legacyscheme.Codecs.UniversalDecoder(), data)
	if err != nil {
		t.Errorf("%v\nData: %s\nSource: %#v", err, string(data), obj)
		return nil
	}
	obj3 := reflect.New(reflect.TypeOf(obj).Elem()).Interface().(runtime.Object)
	err = legacyscheme.Scheme.Convert(obj2, obj3, nil)
	if err != nil {
		t.Errorf("%v\nSource: %#v", err, obj2)
		return nil
	}
	return obj3
}
